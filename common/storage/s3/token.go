package s3

import (
	"fmt"
	"log"
	"time"

	"github.com/2JYK2/go-lib/common/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/dgrijalva/jwt-go"
	"github.com/goccy/go-json"
)

type CustomLogger struct{}

func (c *CustomLogger) Log(args ...interface{}) {
	log.Println(args...) // 将日志输出到标准输出
}

// S3TokenManager S3临时凭证管理器
type S3TokenManager struct {
	AccessKeyId     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	RegionName      string `json:"regionName"`
	BucketName      string `json:"bucketName"`
	PolicyUrl       string `json:"policyUrl"`
	RoleArn         string `json:"roleArn"`
	ServerName      string `json:"serverName"`
	Duration        int    `json:"duration"`
	CDNBaseUrl      string `json:"cdnBaseUrl"`
	CDNSecretKey    string `json:"cdnSecretKey"`
}

// NewS3TokenManager 创建S3临时凭证管理器
func NewS3TokenManager() storage.TokenManager {
	return &S3TokenManager{}
}

// GetProvider 返回存储提供商类型
func (s *S3TokenManager) GetProvider() storage.StorageProvider {
	return storage.ProviderS3
}

// GetSessionToken 保留原有的实现，用于向后兼容
func (s *S3TokenManager) GetSessionToken() (*storage.UploadConfig, error) {
	policyBytes, err := s.buildPolicy()
	if err != nil {
		return nil, err
	}
	// create session
	s1, err := session.NewSession(&aws.Config{
		Region:      &s.RegionName,
		Credentials: credentials.NewStaticCredentials(s.AccessKeyId, s.SecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	svc := sts.New(s1)
	var creds *sts.Credentials
	if s.RoleArn != "" {
		input := &sts.AssumeRoleInput{
			RoleArn:         aws.String(s.RoleArn),
			RoleSessionName: aws.String(s.ServerName),
			Policy:          aws.String(policyBytes),
			DurationSeconds: aws.Int64(int64(s.Duration)),
			Tags:            []*sts.Tag{},
		}

		result, err := svc.AssumeRole(input)
		if err != nil {
			return nil, err
		}
		creds = result.Credentials
	} else {
		input := &sts.GetFederationTokenInput{
			DurationSeconds: aws.Int64(int64(s.Duration)),
			Name:            aws.String(s.ServerName),
			Policy:          aws.String(policyBytes),
			Tags:            []*sts.Tag{},
		}
		result, err := svc.GetFederationToken(input)
		if err != nil {
			return nil, err
		}
		creds = result.Credentials
	}

	return &storage.UploadConfig{
		AccessKeyId:     *creds.AccessKeyId,
		SecretAccessKey: *creds.SecretAccessKey,
		SessionToken:    *creds.SessionToken,
		BucketName:      s.BucketName,
		RegionName:      s.RegionName,
	}, nil
}

// getCOSObjectActions 获取COS对象操作权限
func (s *S3TokenManager) getObjectActions() []string {
	return []string{
		"s3:PutObject",
		"s3:GetObject",
		"s3:DeleteObject",
		"s3:HeadObject",
		"s3:PutObjectAcl",
		"s3:GetObjectAcl",
		"s3:CopyObject",
	}
}

// BuildPolicy 构建COS访问策略
func (s *S3TokenManager) buildPolicy() (string, error) {
	if s.BucketName == "" || s.RegionName == "" {
		return "", fmt.Errorf("bucket, region  are required for COS")
	}
	// 构建COS资源格式
	policy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":   "Allow",
				"Action":   s.getObjectActions(),
				"Resource": s.buildResource(),
			},
		},
	}

	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(policyBytes), nil
}

// getBucketActions 获取存储桶操作权限
func (s *S3TokenManager) getBucketActions() []string {
	return []string{
		"s3:ListBucket",
		"s3:GetBucketLocation",
		"s3:ListBucketMultipartUploads",
	}
}

// buildCOSResource 构建COS资源格式
func (s *S3TokenManager) buildResource() string {
	if s.PolicyUrl != "" {
		return fmt.Sprintf("arn:aws:s3:::%s/%s/*", s.BucketName, s.PolicyUrl)
	}
	return fmt.Sprintf("arn:aws:s3:::%s/*", s.BucketName)
}

type CustomClaims struct {
	jwt.StandardClaims
}

func (s *S3TokenManager) GetCdnUrlAndToken(duration int, serverName, url string) (string, string) {
	claims := CustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(duration)).Unix(), //
			IssuedAt:  time.Now().Unix(),                                            //
			Issuer:    serverName,                                                   //
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.CDNSecretKey))
	if err != nil {
		//logx.Zap().Infow("generateToken SignedString fail", "err", err)
		return "", ""
	}
	return s.CDNBaseUrl, tokenString
}
