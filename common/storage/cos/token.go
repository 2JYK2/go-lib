package cos

import (
	"fmt"
	"log"
	"time"

	"github.com/2JYK2/go-lib/common/storage"
	string2 "github.com/2JYK2/go-lib/common/string"

	sts "github.com/tencentyun/qcloud-cos-sts-sdk/go"
)

type CustomLogger struct{}

func (c *CustomLogger) Log(args ...interface{}) {
	log.Println(args...) // 将日志输出到标准输出
}

// COSTokenManager COS临时凭证管理器
type COSTokenManager struct {
	AccessKeyId     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	RegionName      string `json:"regionName"`
	BucketName      string `json:"bucketName"`
	Uin             string `json:"Uin"`
	PolicyUrl       string `json:"policyUrl"`
	RoleArn         string `json:"roleArn"`
	ServerName      string `json:"serverName"`
	Duration        int    `json:"duration"`
	CDNBaseUrl      string `json:"cdnBaseUrl"`
	CDNSecretKey    string `json:"cdnSecretKey"`
}

// NewCOSTokenManager 创建COS临时凭证管理器
func NewCOSTokenManager() storage.TokenManager {
	return &COSTokenManager{}
}

// GetProvider 返回存储提供商类型
func (c *COSTokenManager) GetProvider() storage.StorageProvider {
	return storage.ProviderCOS
}

// GetSessionTokenLegacy 保留原有的实现，用于向后兼容
func (c *COSTokenManager) GetSessionToken() (*storage.UploadConfig, error) {
	policy, err := c.buildPolicy()
	if err != nil {
		return nil, err
	}
	cliect := sts.NewClient(
		c.AccessKeyId,     // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		c.SecretAccessKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考https://cloud.tencent.com/document/product/598/37140
		nil,
		// sts.Host("sts.internal.tencentcloudapi.com"), // 设置域名, 默认域名sts.tencentcloudapi.com
		// sts.Scheme("http"),      // 设置协议, 默认为https，公有云sts获取临时密钥不允许走http，特殊场景才需要设置http
	)

	// 策略概述 https://cloud.tencent.com/document/product/436/18023
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(c.Duration),
		Region:          c.RegionName,
		Policy:          policy,
		RoleSessionName: c.ServerName,
	}

	var creds *sts.CredentialResult
	if c.RoleArn != "" {
		opt.RoleArn = c.RoleArn
		creds, err = cliect.GetRoleCredential(opt)
		if err != nil {
			panic(err)
		}
	} else {
		creds, err = cliect.GetCredential(opt)
		if err != nil {
			return nil, err
		}
	}

	if creds.Credentials == nil {
		return nil, fmt.Errorf("credentials is nil")
	}

	return &storage.UploadConfig{
		AccessKeyId:     creds.Credentials.TmpSecretID,
		SecretAccessKey: creds.Credentials.TmpSecretKey,
		SessionToken:    creds.Credentials.SessionToken,
		BucketName:      c.BucketName,
		RegionName:      c.RegionName,
	}, nil
}

// BuildPolicy 构建COS访问策略
func (c *COSTokenManager) buildPolicy() (*sts.CredentialPolicy, error) {
	if c.BucketName == "" || c.RegionName == "" || c.Uin == "" {
		return nil, fmt.Errorf("bucket, region and uin are required for COS")
	}
	// 构建COS资源格式
	return &sts.CredentialPolicy{
		Statement: []sts.CredentialPolicyStatement{
			{
				Action:   c.getObjectActions(),
				Effect:   "allow",
				Resource: []string{c.buildResource()},
			},
		},
	}, nil

}

// getCOSObjectActions 获取COS对象操作权限
func (c *COSTokenManager) getObjectActions() []string {
	return []string{
		"cos:PutObject",
		"cos:GetObject",
		"cos:UploadPart",
		"cos:DeleteObject",
		"cos:HeadObject",
		"cos:InitiateMultipartUpload",
		"cos:ListMultipartUploads",
		"cos:ListParts",
		"cos:CompleteMultipartUpload",
	}
}

// buildCOSResource 构建COS资源格式
func (c *COSTokenManager) buildResource() string {
	if c.PolicyUrl != "" {
		return fmt.Sprintf("qcs::cos:%s:uin/%s:%s/%s/*", c.RegionName, c.Uin, c.BucketName, c.PolicyUrl)
	}
	return fmt.Sprintf("qcs::cos:%s:uin/%s:%s/*", c.RegionName, c.Uin, c.BucketName)
}

func (c *COSTokenManager) GetCdnUrlAndToken(duration int, serverName, fileNameUrl string) (string, string) {
	token, err := c.generateCdnToken(c.CDNSecretKey, fileNameUrl, int64(duration))
	if err != nil {
		panic(err)
	}

	return c.CDNBaseUrl, token
}

// 生成 CDN 下载防盗链 Token
func (c *COSTokenManager) generateCdnToken(secretKey, filePath string, expireSeconds int64) (string, error) {
	// 过期时间戳，单位秒
	expireTimestamp := time.Now().Unix() + expireSeconds

	randString := string2.GenerateRandomString(8)
	// 组成签名字符串，格式根据腾讯云文档调整
	// 一般格式： secretKey + filePath + expireTimestamp
	plainText := fmt.Sprintf("%v-%v-%v", expireTimestamp, randString, 0)

	plainTextSign := fmt.Sprintf("%v-%v-%v", filePath, plainText, secretKey)
	// 计算 MD5
	sign := string2.MD5(plainTextSign)

	// 生成带签名的 Token 参数，例如：sign=xxx&t=expireTimestamp
	token := fmt.Sprintf("%v-%v", plainText, sign)
	return token, nil
}
