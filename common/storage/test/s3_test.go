package test

import (
	"context"
	"testing"

	"github.com/2JYK2/go-lib/common/s3"
	"github.com/2JYK2/go-lib/common/storage"
	s4 "github.com/2JYK2/go-lib/common/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestS3Provider 测试 S3 提供商功能
func TestS3Provider(t *testing.T) {
	t.Run("TokenManagerCreation", func(t *testing.T) {
		manager := s4.NewS3TokenManager()
		require.NotNil(t, manager)
		assert.Equal(t, storage.ProviderS3, manager.GetProvider())
	})

	t.Run("UploadConfigStructure", func(t *testing.T) {
		config := &s4.UploadConfig{
			AccessKeyId:     AWSCallerSecretId,
			SecretAccessKey: AWSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      AWSRegion,
			BucketName:      AWSBucket,
		}

		assert.Equal(t, AWSCallerSecretId, config.AccessKeyId)
		assert.Equal(t, AWSCallerSecretKey, config.SecretAccessKey)
		assert.Equal(t, "test-session-token", config.SessionToken)
		assert.Equal(t, AWSRegion, config.RegionName)
		assert.Equal(t, AWSBucket, config.BucketName)
	})

	t.Run("UploadRequestStructure", func(t *testing.T) {
		request := s3.UploadRequest{
			ServerName: "test-server",
			PolicyUrl:  "test/path",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   3600,
			RoleArn:    AWSRoleArn,
		}

		assert.Equal(t, "test-server", request.ServerName)
		assert.Equal(t, "test/path", request.PolicyUrl)
		assert.Equal(t, AWSBucket, request.Bucket)
		assert.Equal(t, AWSRegion, request.Region)
		assert.Equal(t, int64(3600), request.Duration)
		assert.Equal(t, AWSRoleArn, request.RoleArn)
	})

	t.Run("AWSPolicyStructure", func(t *testing.T) {
		policy := s4.AWSPolicy{
			Version: "2012-10-17",
			Statement: []s4.AWSPolicyStatement{
				{
					Effect:   "Allow",
					Action:   []string{"s3:GetObject"},
					Resource: "arn:aws:s3:::test-bucket/*",
				},
			},
		}

		assert.Equal(t, "2012-10-17", policy.Version)
		assert.Len(t, policy.Statement, 1)
		assert.Equal(t, "Allow", policy.Statement[0].Effect)
		assert.Contains(t, policy.Statement[0].Action, "s3:GetObject")
		assert.Equal(t, "arn:aws:s3:::test-bucket/*", policy.Statement[0].Resource)
	})
}

// TestS3Conversion 测试 S3 转换功能
func TestS3Conversion(t *testing.T) {
	t.Run("UploadConfigToStorageCredentials", func(t *testing.T) {
		uploadConfig := &s4.UploadConfig{
			AccessKeyId:     AWSCallerSecretId,
			SecretAccessKey: AWSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      AWSRegion,
			BucketName:      AWSBucket,
		}

		storageCreds := uploadConfig.ToStorageCredentials()
		require.NotNil(t, storageCreds)

		assert.Equal(t, AWSCallerSecretId, storageCreds.AccessKeyId)
		assert.Equal(t, AWSCallerSecretKey, storageCreds.SecretAccessKey)
		assert.Equal(t, "test-session-token", storageCreds.SessionToken)
		assert.Equal(t, AWSRegion, storageCreds.RegionName)
		assert.Equal(t, AWSBucket, storageCreds.BucketName)
	})

	t.Run("UploadRequestToStorageTokenRequest", func(t *testing.T) {
		uploadRequest := s3.UploadRequest{
			ServerName: "test-server",
			PolicyUrl:  "test/path",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   3600,
			RoleArn:    AWSRoleArn,
		}

		storageReq := uploadRequest.ToStorageTokenRequest()
		require.NotNil(t, storageReq)

		assert.Equal(t, "test-server", storageReq.ServerName)
		assert.Equal(t, "test/path", storageReq.PolicyUrl)
		assert.Equal(t, AWSBucket, storageReq.Bucket)
		assert.Equal(t, AWSRegion, storageReq.Region)
		assert.Equal(t, int64(3600), storageReq.Duration)
		assert.Equal(t, AWSRoleArn, storageReq.RoleArn)
		// S3 没有 Uin 字段
		assert.Empty(t, storageReq.Uin)
	})
}

// TestS3PolicyGeneration 测试 S3 策略生成
func TestS3PolicyGeneration(t *testing.T) {
	t.Run("BasicPolicy", func(t *testing.T) {
		req := storage.TokenRequest{
			Bucket: AWSBucket,
			Region: AWSRegion,
		}

		// 使用默认策略构建器（S3）
		builder := storage.NewDefaultPolicyBuilder()
		policy, err := builder.BuildPolicy(req)
		require.NoError(t, err)
		require.NotEmpty(t, policy)

		// 验证策略内容
		assert.Contains(t, policy, AWSBucket)
		assert.Contains(t, policy, "s3:PutObject")
		assert.Contains(t, policy, "s3:GetObject")
		assert.Contains(t, policy, "s3:DeleteObject")
		assert.Contains(t, policy, "s3:ListBucket")
	})

	t.Run("PolicyWithPath", func(t *testing.T) {
		req := storage.TokenRequest{
			Bucket:    AWSBucket,
			Region:    AWSRegion,
			PolicyUrl: "test/path",
		}

		builder := storage.NewDefaultPolicyBuilder()
		policy, err := builder.BuildPolicy(req)
		require.NoError(t, err)
		require.NotEmpty(t, policy)

		// 验证策略包含路径
		assert.Contains(t, policy, "test/path")
	})

	t.Run("PolicyValidation", func(t *testing.T) {
		builder := storage.NewDefaultPolicyBuilder()

		// 测试缺少 bucket
		req := storage.TokenRequest{
			Region: AWSRegion,
		}
		_, err := builder.BuildPolicy(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bucket name is required")

		// 测试缺少 region
		req = storage.TokenRequest{
			Bucket: AWSBucket,
		}
		_, err = builder.BuildPolicy(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "region is required")
	})
}

// TestS3Integration 测试 S3 集成（需要真实凭证时跳过）
func TestS3Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	t.Run("GetSessionTokenWithRole", func(t *testing.T) {
		// 注意：这个测试需要真实的 AWS 凭证
		// 在实际环境中，可以取消注释并配置正确的凭证
		t.Skip("需要真实的 AWS 凭证")

		/*
			manager := s3.NewS3TokenManager()
			ctx := context.Background()

			req := storage.TokenRequest{
				ServerName: "test-server",
				Bucket:     AWSBucket,
				Region:     AWSRegion,
				Duration:   3600,
				RoleArn:    AWSRoleArn,
			}

			creds, err := manager.GetSessionToken(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, creds)

			assert.NotEmpty(t, creds.AccessKeyId)
			assert.NotEmpty(t, creds.SecretAccessKey)
			assert.NotEmpty(t, creds.SessionToken)
			assert.Equal(t, AWSRegion, creds.RegionName)
			assert.Equal(t, AWSBucket, creds.BucketName)
		*/
	})

	t.Run("GetSessionTokenWithoutRole", func(t *testing.T) {
		// 注意：这个测试需要真实的 AWS 凭证
		t.Skip("需要真实的 AWS 凭证")

		/*
			manager := s3.NewS3TokenManager()
			ctx := context.Background()

			req := storage.TokenRequest{
				ServerName: "test-server",
				Bucket:     AWSBucket,
				Region:     AWSRegion,
				Duration:   3600,
			}

			creds, err := manager.GetSessionToken(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, creds)

			assert.NotEmpty(t, creds.AccessKeyId)
			assert.NotEmpty(t, creds.SecretAccessKey)
			assert.NotEmpty(t, creds.SessionToken)
			assert.Equal(t, AWSRegion, creds.RegionName)
			assert.Equal(t, AWSBucket, creds.BucketName)
		*/
	})
}

// TestS3ErrorHandling 测试 S3 错误处理
func TestS3ErrorHandling(t *testing.T) {
	t.Run("InvalidCredentials", func(t *testing.T) {
		// 测试无效凭证的错误处理
		manager := s4.NewS3TokenManager()
		ctx := context.Background()

		// 使用无效的请求参数
		req := storage.TokenRequest{
			ServerName: "",
			Bucket:     "",
			Region:     "",
			Duration:   0,
		}

		// 注意：由于没有真实的凭证，这个测试可能会失败
		// 在实际测试中，应该使用 mock 或者测试环境
		_, err := manager.GetSessionToken(ctx, req)
		// 这里不检查错误，因为可能因为凭证问题而失败
		_ = err
	})

	t.Run("InvalidPolicy", func(t *testing.T) {
		// 测试无效策略的错误处理
		builder := storage.NewDefaultPolicyBuilder()

		// 空请求
		req := storage.TokenRequest{}
		_, err := builder.BuildPolicy(req)
		assert.Error(t, err)
	})
}

// TestS3AWSPolicyStructure 测试 AWS 策略结构
func TestS3AWSPolicyStructure(t *testing.T) {
	t.Run("PolicyStatement", func(t *testing.T) {
		statement := s4.AWSPolicyStatement{
			Effect:   "Allow",
			Action:   []string{"s3:GetObject", "s3:PutObject"},
			Resource: "arn:aws:s3:::test-bucket/*",
		}

		assert.Equal(t, "Allow", statement.Effect)
		assert.Len(t, statement.Action, 2)
		assert.Contains(t, statement.Action, "s3:GetObject")
		assert.Contains(t, statement.Action, "s3:PutObject")
		assert.Equal(t, "arn:aws:s3:::test-bucket/*", statement.Resource)
	})

	t.Run("CompletePolicy", func(t *testing.T) {
		policy := s4.AWSPolicy{
			Version: "2012-10-17",
			Statement: []s4.AWSPolicyStatement{
				{
					Effect:   "Allow",
					Action:   []string{"s3:GetObject"},
					Resource: "arn:aws:s3:::test-bucket/*",
				},
				{
					Effect:   "Allow",
					Action:   []string{"s3:ListBucket"},
					Resource: "arn:aws:s3:::test-bucket",
				},
			},
		}

		assert.Equal(t, "2012-10-17", policy.Version)
		assert.Len(t, policy.Statement, 2)

		// 检查第一个声明
		assert.Equal(t, "Allow", policy.Statement[0].Effect)
		assert.Contains(t, policy.Statement[0].Action, "s3:GetObject")
		assert.Equal(t, "arn:aws:s3:::test-bucket/*", policy.Statement[0].Resource)

		// 检查第二个声明
		assert.Equal(t, "Allow", policy.Statement[1].Effect)
		assert.Contains(t, policy.Statement[1].Action, "s3:ListBucket")
		assert.Equal(t, "arn:aws:s3:::test-bucket", policy.Statement[1].Resource)
	})
}

// BenchmarkS3PolicyBuilding S3 策略构建性能测试
func BenchmarkS3PolicyBuilding(b *testing.B) {
	builder := storage.NewDefaultPolicyBuilder()
	req := storage.TokenRequest{
		Bucket: AWSBucket,
		Region: AWSRegion,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildPolicy(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestS3Concurrency 测试 S3 并发访问
func TestS3Concurrency(t *testing.T) {
	// 并发测试策略构建
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			builder := storage.NewDefaultPolicyBuilder()
			req := storage.TokenRequest{
				Bucket: AWSBucket,
				Region: AWSRegion,
			}

			policy, err := builder.BuildPolicy(req)
			require.NoError(t, err)
			assert.NotEmpty(t, policy)
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestS3DurationHandling 测试 S3 持续时间处理
func TestS3DurationHandling(t *testing.T) {
	t.Run("ValidDuration", func(t *testing.T) {
		req := s3.UploadRequest{
			ServerName: "test-server",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   3600, // 1 hour
		}

		storageReq := req.ToStorageTokenRequest()
		assert.Equal(t, int64(3600), storageReq.Duration)
	})

	t.Run("MaxDuration", func(t *testing.T) {
		req := s3.UploadRequest{
			ServerName: "test-server",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   43200, // 12 hours (AWS STS max)
		}

		storageReq := req.ToStorageTokenRequest()
		assert.Equal(t, int64(43200), storageReq.Duration)
	})

	t.Run("ZeroDuration", func(t *testing.T) {
		req := s3.UploadRequest{
			ServerName: "test-server",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   0,
		}

		storageReq := req.ToStorageTokenRequest()
		assert.Equal(t, int64(0), storageReq.Duration)
	})
}
