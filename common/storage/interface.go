package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// StorageProvider 定义对象存储提供商的类型
type StorageProvider int

const (
	ProviderCOS StorageProvider = 1
	ProviderS3  StorageProvider = 0
)

// Credentials 定义存储服务的认证信息
type UploadConfig struct {
	AccessKeyId     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	RegionName      string `json:"regionName"`
	BucketName      string `json:"bucketName"`
}

type DeleteResponse struct {
	// Container element for a successful delete. It identifies the object that was
	// successfully deleted.
	Deleted []types.DeletedObject

	// Container for a failed delete action that describes the object that Amazon S3
	// attempted to delete and the error it encountered.
	Errors []types.Error
}

// TokenManager 定义临时凭证管理接口
type TokenManager interface {
	// GetProvider 返回存储提供商类型
	GetProvider() StorageProvider

	// GetSessionToken 获取临时访问凭证
	GetSessionToken() (*UploadConfig, error)

	// GetCdnUrlAndToken 获取cdn的url 和 token
	GetCdnUrlAndToken(duration int, serverName, url string) (string, string)

	// GetClient 获取一个客户端，然后做操作
	InitClient() error

	// 批量删除
	BatchDelete(ctx context.Context, keys []string) (*DeleteResponse, error)
}
