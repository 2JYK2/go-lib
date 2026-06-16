package test

import (
	"testing"

	"github.com/2JYK2/go-lib/common/cos"
	"github.com/2JYK2/go-lib/common/s3"
	"github.com/2JYK2/go-lib/common/storage"
	cos2 "github.com/2JYK2/go-lib/common/storage/cos"
	s4 "github.com/2JYK2/go-lib/common/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageProviderTypes 测试存储提供商类型
func TestStorageProviderTypes(t *testing.T) {
	t.Run("ProviderTypes", func(t *testing.T) {
		// 测试存储提供商类型定义
		assert.Equal(t, "cos", string(storage.ProviderCOS))
		assert.Equal(t, "s3", string(storage.ProviderS3))
	})

	t.Run("CredentialsStructure", func(t *testing.T) {
		// 测试凭证结构
		creds := &storage.Credentials{
			AccessKeyId:     "test-access-key",
			SecretAccessKey: "test-secret-key",
			SessionToken:    "test-session-token",
			RegionName:      "us-east-1",
			BucketName:      "test-bucket",
		}

		assert.NotEmpty(t, creds.AccessKeyId)
		assert.NotEmpty(t, creds.SecretAccessKey)
		assert.NotEmpty(t, creds.SessionToken)
		assert.NotEmpty(t, creds.RegionName)
		assert.NotEmpty(t, creds.BucketName)
	})

	t.Run("TokenRequestStructure", func(t *testing.T) {
		// 测试临时凭证请求结构
		req := storage.TokenRequest{
			ServerName: "test-server",
			PolicyUrl:  "test/path",
			Bucket:     "test-bucket",
			Region:     "us-east-1",
			Duration:   3600,
			RoleArn:    "arn:aws:iam::123456789012:role/test-role",
			Uin:        "123456789",
		}

		assert.Equal(t, "test-server", req.ServerName)
		assert.Equal(t, "test/path", req.PolicyUrl)
		assert.Equal(t, "test-bucket", req.Bucket)
		assert.Equal(t, "us-east-1", req.Region)
		assert.Equal(t, int64(3600), req.Duration)
		assert.Equal(t, "arn:aws:iam::123456789012:role/test-role", req.RoleArn)
		assert.Equal(t, "123456789", req.Uin)
	})
}

// TestTokenManagerFactory 测试临时凭证管理器工厂
func TestTokenManagerFactory(t *testing.T) {
	factory := storage.NewDefaultTokenManagerFactory()

	t.Run("RegisterAndCreateCOSManager", func(t *testing.T) {
		// 注册 COS 管理器
		factory.RegisterManager(storage.ProviderCOS, func() (storage.TokenManager, error) {
			return cos2.NewCOSTokenManager(), nil
		})

		// 创建 COS 管理器
		manager, err := factory.CreateTokenManager(storage.ProviderCOS)
		require.NoError(t, err)
		require.NotNil(t, manager)
		assert.Equal(t, storage.ProviderCOS, manager.GetProvider())
	})

	t.Run("RegisterAndCreateS3Manager", func(t *testing.T) {
		// 注册 S3 管理器
		factory.RegisterManager(storage.ProviderS3, func() (storage.TokenManager, error) {
			return s4.NewS3TokenManager(), nil
		})

		// 创建 S3 管理器
		manager, err := factory.CreateTokenManager(storage.ProviderS3)
		require.NoError(t, err)
		require.NotNil(t, manager)
		assert.Equal(t, storage.ProviderS3, manager.GetProvider())
	})

	t.Run("UnsupportedProvider", func(t *testing.T) {
		// 测试不支持的提供商
		manager, err := factory.CreateTokenManager("unsupported")
		assert.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "unsupported storage provider")
	})
}

// TestTokenManagerRegistry 测试临时凭证管理器注册表
func TestTokenManagerRegistry(t *testing.T) {
	factory := storage.NewDefaultTokenManagerFactory()
	registry := storage.NewTokenManagerRegistry(factory)

	// 注册管理器
	factory.RegisterManager(storage.ProviderCOS, func() (storage.TokenManager, error) {
		return cos2.NewCOSTokenManager(), nil
	})
	factory.RegisterManager(storage.ProviderS3, func() (storage.TokenManager, error) {
		return s4.NewS3TokenManager(), nil
	})

	t.Run("GetSupportedProviders", func(t *testing.T) {
		providers := registry.GetSupportedProviders()
		assert.Contains(t, providers, storage.ProviderCOS)
		assert.Contains(t, providers, storage.ProviderS3)
		assert.Len(t, providers, 2)
	})

	t.Run("GetTokenManager", func(t *testing.T) {
		// 获取 COS 管理器
		cosManager, err := registry.GetTokenManager(storage.ProviderCOS)
		require.NoError(t, err)
		assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())

		// 获取 S3 管理器
		s3Manager, err := registry.GetTokenManager(storage.ProviderS3)
		require.NoError(t, err)
		assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
	})
}

// TestPolicyBuilder 测试策略构建器
func TestPolicyBuilder(t *testing.T) {
	t.Run("DefaultPolicyBuilder", func(t *testing.T) {
		builder := storage.NewDefaultPolicyBuilder()

		req := storage.TokenRequest{
			Bucket: "test-bucket",
			Region: "us-east-1",
		}

		policy, err := builder.BuildPolicy(req)
		require.NoError(t, err)
		assert.NotEmpty(t, policy)
		assert.Contains(t, policy, "test-bucket")
		assert.Contains(t, policy, "s3:PutObject")
		assert.Contains(t, policy, "s3:GetObject")
	})

	t.Run("COSPolicyBuilder", func(t *testing.T) {
		builder := storage.NewCOSPolicyBuilder()

		req := storage.TokenRequest{
			Bucket: "test-bucket",
			Region: "ap-tokyo",
			Uin:    "123456789",
		}

		policy, err := builder.BuildPolicy(req)
		require.NoError(t, err)
		assert.NotEmpty(t, policy)
		assert.Contains(t, policy, "test-bucket")
		assert.Contains(t, policy, "ap-tokyo")
		assert.Contains(t, policy, "123456789")
		assert.Contains(t, policy, "cos:PutObject")
		assert.Contains(t, policy, "cos:GetObject")
	})

	t.Run("GetPolicyBuilder", func(t *testing.T) {
		// 测试根据提供商获取策略构建器
		cosBuilder := storage.GetPolicyBuilder(storage.ProviderCOS)
		assert.IsType(t, &storage.COSPolicyBuilder{}, cosBuilder)

		s3Builder := storage.GetPolicyBuilder(storage.ProviderS3)
		assert.IsType(t, &storage.DefaultPolicyBuilder{}, s3Builder)

		defaultBuilder := storage.GetPolicyBuilder("unknown")
		assert.IsType(t, &storage.DefaultPolicyBuilder{}, defaultBuilder)
	})
}

// TestCOSCompatibility 测试 COS 兼容性
func TestCOSCompatibility(t *testing.T) {
	t.Run("UploadConfigConversion", func(t *testing.T) {
		// 测试 COS 上传配置转换
		uploadConfig := &cos.UploadConfig{
			AccessKeyId:     COSCallerSecretId,
			SecretAccessKey: COSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      COSRegion,
			BucketName:      COSBucket,
		}

		storageCreds := uploadConfig.ToStorageCredentials()
		assert.Equal(t, COSCallerSecretId, storageCreds.AccessKeyId)
		assert.Equal(t, COSCallerSecretKey, storageCreds.SecretAccessKey)
		assert.Equal(t, "test-session-token", storageCreds.SessionToken)
		assert.Equal(t, COSRegion, storageCreds.RegionName)
		assert.Equal(t, COSBucket, storageCreds.BucketName)
	})

	t.Run("UploadRequestConversion", func(t *testing.T) {
		// 测试 COS 上传请求转换
		uploadRequest := cos2.UploadRequest{
			AccessKeyId:     COSCallerSecretId,
			SecretAccessKey: COSCallerSecretKey,
			ServerName:      "test-server",
			PolicyUrl:       "test/path",
			Bucket:          COSBucket,
			Region:          COSRegion,
			Duration:        3600,
			RoleArn:         COSRoleArn,
			Uin:             COSUin,
		}

		storageReq := uploadRequest.ToStorageTokenRequest()
		assert.Equal(t, "test-server", storageReq.ServerName)
		assert.Equal(t, "test/path", storageReq.PolicyUrl)
		assert.Equal(t, COSBucket, storageReq.Bucket)
		assert.Equal(t, COSRegion, storageReq.Region)
		assert.Equal(t, int64(3600), storageReq.Duration)
		assert.Equal(t, COSRoleArn, storageReq.RoleArn)
		assert.Equal(t, COSUin, storageReq.Uin)
	})

	t.Run("COSTokenManager", func(t *testing.T) {
		// 测试 COS 临时凭证管理器
		manager := cos2.NewCOSTokenManager()
		assert.Equal(t, storage.ProviderCOS, manager.GetProvider())

		// 注意：这里不实际调用 GetSessionToken，因为需要真实的 AWS 凭证
		// 在实际测试中，可以使用 mock 或者集成测试环境
	})
}

// TestS3Compatibility 测试 S3 兼容性
func TestS3Compatibility(t *testing.T) {
	t.Run("UploadConfigConversion", func(t *testing.T) {
		// 测试 S3 上传配置转换
		uploadConfig := &s4.UploadConfig{
			AccessKeyId:     AWSCallerSecretId,
			SecretAccessKey: AWSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      AWSRegion,
			BucketName:      AWSBucket,
		}

		storageCreds := uploadConfig.ToStorageCredentials()
		assert.Equal(t, AWSCallerSecretId, storageCreds.AccessKeyId)
		assert.Equal(t, AWSCallerSecretKey, storageCreds.SecretAccessKey)
		assert.Equal(t, "test-session-token", storageCreds.SessionToken)
		assert.Equal(t, AWSRegion, storageCreds.RegionName)
		assert.Equal(t, AWSBucket, storageCreds.BucketName)
	})

	t.Run("UploadRequestConversion", func(t *testing.T) {
		// 测试 S3 上传请求转换
		uploadRequest := s3.UploadRequest{
			ServerName: "test-server",
			PolicyUrl:  "test/path",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   3600,
			RoleArn:    AWSRoleArn,
		}

		storageReq := uploadRequest.ToStorageTokenRequest()
		assert.Equal(t, "test-server", storageReq.ServerName)
		assert.Equal(t, "test/path", storageReq.PolicyUrl)
		assert.Equal(t, AWSBucket, storageReq.Bucket)
		assert.Equal(t, AWSRegion, storageReq.Region)
		assert.Equal(t, int64(3600), storageReq.Duration)
		assert.Equal(t, AWSRoleArn, storageReq.RoleArn)
	})

	t.Run("S3TokenManager", func(t *testing.T) {
		// 测试 S3 临时凭证管理器
		manager := s4.NewS3TokenManager()
		assert.Equal(t, storage.ProviderS3, manager.GetProvider())

		// 注意：这里不实际调用 GetSessionToken，因为需要真实的 AWS 凭证
		// 在实际测试中，可以使用 mock 或者集成测试环境
	})
}

// TestIntegration 集成测试（需要真实的凭证，通常跳过）
func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	t.Run("COSIntegration", func(t *testing.T) {
		// 这里可以添加真实的 COS 集成测试
		// 需要真实的腾讯云凭证
		t.Skip("需要真实的腾讯云凭证")
	})

	t.Run("S3Integration", func(t *testing.T) {
		// 这里可以添加真实的 S3 集成测试
		// 需要真实的 AWS 凭证
		t.Skip("需要真实的 AWS 凭证")
	})
}

// TestErrorHandling 测试错误处理
func TestErrorHandling(t *testing.T) {
	t.Run("PolicyBuilderErrors", func(t *testing.T) {
		builder := storage.NewDefaultPolicyBuilder()

		// 测试缺少 bucket 的错误
		req := storage.TokenRequest{
			Region: "us-east-1",
		}
		_, err := builder.BuildPolicy(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bucket name is required")

		// 测试缺少 region 的错误
		req = storage.TokenRequest{
			Bucket: "test-bucket",
		}
		_, err = builder.BuildPolicy(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "region is required")
	})

	t.Run("COSPolicyBuilderErrors", func(t *testing.T) {
		builder := storage.NewCOSPolicyBuilder()

		// 测试缺少必需字段的错误
		req := storage.TokenRequest{
			Bucket: "test-bucket",
		}
		_, err := builder.BuildPolicy(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bucket, region and uin are required for COS")
	})
}

// BenchmarkPolicyBuilding 性能测试
func BenchmarkPolicyBuilding(b *testing.B) {
	builder := storage.NewDefaultPolicyBuilder()
	req := storage.TokenRequest{
		Bucket: "test-bucket",
		Region: "us-east-1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildPolicy(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestConcurrentAccess 并发访问测试
func TestConcurrentAccess(t *testing.T) {
	factory := storage.NewDefaultTokenManagerFactory()
	factory.RegisterManager(storage.ProviderCOS, func() (storage.TokenManager, error) {
		return cos2.NewCOSTokenManager(), nil
	})
	factory.RegisterManager(storage.ProviderS3, func() (storage.TokenManager, error) {
		return s4.NewS3TokenManager(), nil
	})

	// 并发创建管理器
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			cosManager, err := factory.CreateTokenManager(storage.ProviderCOS)
			require.NoError(t, err)
			assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())

			s3Manager, err := factory.CreateTokenManager(storage.ProviderS3)
			require.NoError(t, err)
			assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}
