package s3

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Clients = make(map[string]*s3.Client)
	mu        sync.Mutex
)

// GetClient 获取或创建 S3 Client（按 bucket 缓存）
func (s *S3TokenManager) InitClient() error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := s3Clients[s.BucketName]; ok {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(s.RegionName),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				s.AccessKeyId,
				s.SecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UseAccelerate = false
	})

	s3Clients[s.BucketName] = client
	return nil
}

func (s *S3TokenManager) getClient() (*s3.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if c, ok := s3Clients[s.BucketName]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("not found s3 client %s", s.BucketName)
}
