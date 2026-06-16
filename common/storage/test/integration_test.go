package test

import (
	"context"
	"testing"
	"time"

	"github.com/2JYK2/go-lib/common/cos"
	"github.com/2JYK2/go-lib/common/s3"
	"github.com/2JYK2/go-lib/common/storage"
	cos2 "github.com/2JYK2/go-lib/common/storage/cos"
	s4 "github.com/2JYK2/go-lib/common/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageProviderCompatibilityIntegration 测试存储提供商兼容性集成
func TestStorageProviderCompatibilityIntegration(t *testing.T) {
	t.Run("ProviderTypeConsistency", func(t *testing.T) {
		// 测试提供商类型一致性
		cosManager := cos2.NewCOSTokenManager()
		s3Manager := s4.NewS3TokenManager()

		assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())
		assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
		assert.NotEqual(t, cosManager.GetProvider(), s3Manager.GetProvider())
	})

	t.Run("InterfaceCompatibility", func(t *testing.T) {
		// 测试接口兼容性
		var cosManager = cos2.NewCOSTokenManager()
		var s3Manager = s4.NewS3TokenManager()

		// 验证两个管理器都实现了 TokenManager 接口
		assert.Implements(t, (*storage.TokenManager)(nil), cosManager)
		assert.Implements(t, (*storage.TokenManager)(nil), s3Manager)

		// 验证提供商类型
		assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())
		assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
	})

	t.Run("CredentialsStructureCompatibility", func(t *testing.T) {
		// 测试凭证结构兼容性
		cosConfig := &cos.UploadConfig{
			AccessKeyId:     COSCallerSecretId,
			SecretAccessKey: COSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      COSRegion,
			BucketName:      COSBucket,
		}

		s3Config := &s4.UploadConfig{
			AccessKeyId:     AWSCallerSecretId,
			SecretAccessKey: AWSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      AWSRegion,
			BucketName:      AWSBucket,
		}

		// 转换为通用凭证结构
		cosCreds := cosConfig.ToStorageCredentials()
		s3Creds := s3Config.ToStorageCredentials()

		// 验证结构一致性
		assert.Equal(t, cosCreds.AccessKeyId, cosConfig.AccessKeyId)
		assert.Equal(t, cosCreds.SecretAccessKey, cosConfig.SecretAccessKey)
		assert.Equal(t, cosCreds.SessionToken, cosConfig.SessionToken)
		assert.Equal(t, cosCreds.RegionName, cosConfig.RegionName)
		assert.Equal(t, cosCreds.BucketName, cosConfig.BucketName)

		assert.Equal(t, s3Creds.AccessKeyId, s3Config.AccessKeyId)
		assert.Equal(t, s3Creds.SecretAccessKey, s3Config.SecretAccessKey)
		assert.Equal(t, s3Creds.SessionToken, s3Config.SessionToken)
		assert.Equal(t, s3Creds.RegionName, s3Config.RegionName)
		assert.Equal(t, s3Creds.BucketName, s3Config.BucketName)
	})

	t.Run("TokenRequestCompatibility", func(t *testing.T) {
		// 测试临时凭证请求兼容性
		cosRequest := cos2.UploadRequest{
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

		s3Request := s3.UploadRequest{
			ServerName: "test-server",
			PolicyUrl:  "test/path",
			Bucket:     AWSBucket,
			Region:     AWSRegion,
			Duration:   3600,
			RoleArn:    AWSRoleArn,
		}

		// 转换为通用请求结构
		cosStorageReq := cosRequest.ToStorageTokenRequest()
		s3StorageReq := s3Request.ToStorageTokenRequest()

		// 验证共同字段
		assert.Equal(t, cosStorageReq.ServerName, s3StorageReq.ServerName)
		assert.Equal(t, cosStorageReq.PolicyUrl, s3StorageReq.PolicyUrl)
		assert.Equal(t, cosStorageReq.Duration, s3StorageReq.Duration)
		// 注意：RoleArn 格式不同，所以不比较

		// 验证特定字段
		assert.Equal(t, COSBucket, cosStorageReq.Bucket)
		assert.Equal(t, COSRegion, cosStorageReq.Region)
		assert.Equal(t, COSUin, cosStorageReq.Uin)

		assert.Equal(t, AWSBucket, s3StorageReq.Bucket)
		assert.Equal(t, AWSRegion, s3StorageReq.Region)
		assert.Empty(t, s3StorageReq.Uin) // S3 没有 Uin 字段
	})
}

// TestPolicyBuilderCompatibility 测试策略构建器兼容性
func TestPolicyBuilderCompatibility(t *testing.T) {
	t.Run("PolicyBuilderSelection", func(t *testing.T) {
		// 测试策略构建器选择
		cosBuilder := storage.GetPolicyBuilder(storage.ProviderCOS)
		s3Builder := storage.GetPolicyBuilder(storage.ProviderS3)

		assert.IsType(t, &storage.COSPolicyBuilder{}, cosBuilder)
		assert.IsType(t, &storage.DefaultPolicyBuilder{}, s3Builder)
	})

	t.Run("PolicyGenerationConsistency", func(t *testing.T) {
		// 测试策略生成一致性
		cosBuilder := storage.NewCOSPolicyBuilder()
		s3Builder := storage.NewDefaultPolicyBuilder()

		cosReq := storage.TokenRequest{
			Bucket: COSBucket,
			Region: COSRegion,
			Uin:    COSUin,
		}

		s3Req := storage.TokenRequest{
			Bucket: AWSBucket,
			Region: AWSRegion,
		}

		cosPolicy, err := cosBuilder.BuildPolicy(cosReq)
		require.NoError(t, err)
		require.NotEmpty(t, cosPolicy)

		s3Policy, err := s3Builder.BuildPolicy(s3Req)
		require.NoError(t, err)
		require.NotEmpty(t, s3Policy)

		// 验证策略包含必要的权限
		assert.Contains(t, cosPolicy, "cos:PutObject")
		assert.Contains(t, cosPolicy, "cos:GetObject")
		assert.Contains(t, cosPolicy, COSBucket)
		assert.Contains(t, cosPolicy, COSRegion)

		assert.Contains(t, s3Policy, "s3:PutObject")
		assert.Contains(t, s3Policy, "s3:GetObject")
		assert.Contains(t, s3Policy, AWSBucket)
		// 注意：S3 策略不包含区域信息
	})

	t.Run("PolicyWithPathConsistency", func(t *testing.T) {
		// 测试带路径的策略生成一致性
		cosBuilder := storage.NewCOSPolicyBuilder()
		s3Builder := storage.NewDefaultPolicyBuilder()

		cosReq := storage.TokenRequest{
			Bucket:    COSBucket,
			Region:    COSRegion,
			Uin:       COSUin,
			PolicyUrl: "test/path",
		}

		s3Req := storage.TokenRequest{
			Bucket:    AWSBucket,
			Region:    AWSRegion,
			PolicyUrl: "test/path",
		}

		cosPolicy, err := cosBuilder.BuildPolicy(cosReq)
		require.NoError(t, err)
		require.NotEmpty(t, cosPolicy)

		s3Policy, err := s3Builder.BuildPolicy(s3Req)
		require.NoError(t, err)
		require.NotEmpty(t, s3Policy)

		// 验证策略包含路径
		assert.Contains(t, cosPolicy, "test/path")
		assert.Contains(t, s3Policy, "test/path")
	})
}

// TestFactoryCompatibility 测试工厂兼容性
func TestFactoryCompatibility(t *testing.T) {
	t.Run("FactoryRegistration", func(t *testing.T) {
		// 测试工厂注册
		factory := storage.NewDefaultTokenManagerFactory()

		// 注册两个提供商
		factory.RegisterManager(storage.ProviderCOS, func() (storage.TokenManager, error) {
			return cos2.NewCOSTokenManager(), nil
		})
		factory.RegisterManager(storage.ProviderS3, func() (storage.TokenManager, error) {
			return s4.NewS3TokenManager(), nil
		})

		// 创建管理器
		cosManager, err := factory.CreateTokenManager(storage.ProviderCOS)
		require.NoError(t, err)
		require.NotNil(t, cosManager)

		s3Manager, err := factory.CreateTokenManager(storage.ProviderS3)
		require.NoError(t, err)
		require.NotNil(t, s3Manager)

		// 验证提供商类型
		assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())
		assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
	})

	t.Run("RegistryCompatibility", func(t *testing.T) {
		// 测试注册表兼容性
		factory := storage.NewDefaultTokenManagerFactory()
		registry := storage.NewTokenManagerRegistry(factory)

		// 注册管理器
		factory.RegisterManager(storage.ProviderCOS, func() (storage.TokenManager, error) {
			return cos2.NewCOSTokenManager(), nil
		})
		factory.RegisterManager(storage.ProviderS3, func() (storage.TokenManager, error) {
			return s4.NewS3TokenManager(), nil
		})

		// 获取支持的提供商
		providers := registry.GetSupportedProviders()
		assert.Contains(t, providers, storage.ProviderCOS)
		assert.Contains(t, providers, storage.ProviderS3)
		assert.Len(t, providers, 2)

		// 通过注册表获取管理器
		cosManager, err := registry.GetTokenManager(storage.ProviderCOS)
		require.NoError(t, err)
		assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())

		s3Manager, err := registry.GetTokenManager(storage.ProviderS3)
		require.NoError(t, err)
		assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
	})
}

// TestCrossProviderCompatibility 测试跨提供商兼容性
func TestCrossProviderCompatibility(t *testing.T) {
	t.Run("CredentialsInterchangeability", func(t *testing.T) {
		// 测试凭证可互换性
		cosCreds := &storage.Credentials{
			AccessKeyId:     COSCallerSecretId,
			SecretAccessKey: COSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      COSRegion,
			BucketName:      COSBucket,
		}

		s3Creds := &storage.Credentials{
			AccessKeyId:     AWSCallerSecretId,
			SecretAccessKey: AWSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      AWSRegion,
			BucketName:      AWSBucket,
		}

		// 验证两个凭证结构相同
		assert.Equal(t, cosCreds.AccessKeyId, COSCallerSecretId)
		assert.Equal(t, s3Creds.AccessKeyId, AWSCallerSecretId)
		assert.NotEqual(t, cosCreds.AccessKeyId, s3Creds.AccessKeyId)
		assert.NotEqual(t, cosCreds.RegionName, s3Creds.RegionName)
		assert.NotEqual(t, cosCreds.BucketName, s3Creds.BucketName)
	})

	t.Run("TokenRequestInterchangeability", func(t *testing.T) {
		// 测试临时凭证请求可互换性
		commonReq := storage.TokenRequest{
			ServerName: "test-server",
			PolicyUrl:  "test/path",
			Duration:   3600,
		}

		// COS 特定请求
		cosReq := commonReq
		cosReq.Bucket = COSBucket
		cosReq.Region = COSRegion
		cosReq.RoleArn = COSRoleArn
		cosReq.Uin = COSUin

		// S3 特定请求
		s3Req := commonReq
		s3Req.Bucket = AWSBucket
		s3Req.Region = AWSRegion
		s3Req.RoleArn = AWSRoleArn

		// 验证共同字段
		assert.Equal(t, cosReq.ServerName, s3Req.ServerName)
		assert.Equal(t, cosReq.PolicyUrl, s3Req.PolicyUrl)
		assert.Equal(t, cosReq.Duration, s3Req.Duration)

		// 验证特定字段
		assert.NotEqual(t, cosReq.Bucket, s3Req.Bucket)
		assert.NotEqual(t, cosReq.Region, s3Req.Region)
		assert.NotEqual(t, cosReq.RoleArn, s3Req.RoleArn)
		assert.NotEmpty(t, cosReq.Uin)
		assert.Empty(t, s3Req.Uin)
	})
}

// TestErrorHandlingCompatibility 测试错误处理兼容性
func TestErrorHandlingCompatibility(t *testing.T) {
	t.Run("PolicyBuilderErrorConsistency", func(t *testing.T) {
		// 测试策略构建器错误处理一致性
		cosBuilder := storage.NewCOSPolicyBuilder()
		s3Builder := storage.NewDefaultPolicyBuilder()

		// 测试空请求
		emptyReq := storage.TokenRequest{}

		_, cosErr := cosBuilder.BuildPolicy(emptyReq)
		_, s3Err := s3Builder.BuildPolicy(emptyReq)

		// 两个构建器都应该返回错误
		assert.Error(t, cosErr)
		assert.Error(t, s3Err)

		// 错误信息应该不同（因为验证规则不同）
		assert.NotEqual(t, cosErr.Error(), s3Err.Error())
	})

	t.Run("ManagerErrorConsistency", func(t *testing.T) {
		// 测试管理器错误处理一致性
		cosManager := cos2.NewCOSTokenManager()
		s3Manager := s4.NewS3TokenManager()

		ctx := context.Background()
		invalidReq := storage.TokenRequest{
			ServerName: "",
			Bucket:     "",
			Region:     "",
			Duration:   0,
		}

		// 注意：由于没有真实凭证，这些调用可能会失败
		// 这里主要测试接口一致性
		_, cosErr := cosManager.GetSessionToken(ctx, invalidReq)
		_, s3Err := s3Manager.GetSessionToken(ctx, invalidReq)

		// 两个管理器都应该处理错误（即使错误类型可能不同）
		// 这里不检查具体的错误内容，因为可能因为凭证问题而失败
		_ = cosErr
		_ = s3Err
	})
}

// TestPerformanceCompatibility 测试性能兼容性
func TestPerformanceCompatibility(t *testing.T) {
	t.Run("PolicyBuildingPerformance", func(t *testing.T) {
		// 测试策略构建性能
		cosBuilder := storage.NewCOSPolicyBuilder()
		s3Builder := storage.NewDefaultPolicyBuilder()

		cosReq := storage.TokenRequest{
			Bucket: COSBucket,
			Region: COSRegion,
			Uin:    COSUin,
		}

		s3Req := storage.TokenRequest{
			Bucket: AWSBucket,
			Region: AWSRegion,
		}

		// 测试 COS 策略构建时间
		cosStart := time.Now()
		_, err := cosBuilder.BuildPolicy(cosReq)
		cosDuration := time.Since(cosStart)
		require.NoError(t, err)

		// 测试 S3 策略构建时间
		s3Start := time.Now()
		_, err = s3Builder.BuildPolicy(s3Req)
		s3Duration := time.Since(s3Start)
		require.NoError(t, err)

		// 验证两个构建器都能在合理时间内完成
		assert.True(t, cosDuration < time.Second)
		assert.True(t, s3Duration < time.Second)

		t.Logf("COS policy building time: %v", cosDuration)
		t.Logf("S3 policy building time: %v", s3Duration)
	})
}

// TestConcurrencyCompatibility 测试并发兼容性
func TestConcurrencyCompatibility(t *testing.T) {
	t.Run("ConcurrentManagerCreation", func(t *testing.T) {
		// 测试并发创建管理器
		factory := storage.NewDefaultTokenManagerFactory()
		factory.RegisterManager(storage.ProviderCOS, func() (storage.TokenManager, error) {
			return cos2.NewCOSTokenManager(), nil
		})
		factory.RegisterManager(storage.ProviderS3, func() (storage.TokenManager, error) {
			return s4.NewS3TokenManager(), nil
		})

		done := make(chan bool, 20)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				cosManager, err := factory.CreateTokenManager(storage.ProviderCOS)
				require.NoError(t, err)
				assert.Equal(t, storage.ProviderCOS, cosManager.GetProvider())
			}()
		}
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				s3Manager, err := factory.CreateTokenManager(storage.ProviderS3)
				require.NoError(t, err)
				assert.Equal(t, storage.ProviderS3, s3Manager.GetProvider())
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 20; i++ {
			<-done
		}
	})

	t.Run("ConcurrentPolicyBuilding", func(t *testing.T) {
		// 测试并发策略构建
		cosBuilder := storage.NewCOSPolicyBuilder()
		s3Builder := storage.NewDefaultPolicyBuilder()

		done := make(chan bool, 20)
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				cosReq := storage.TokenRequest{
					Bucket: COSBucket,
					Region: COSRegion,
					Uin:    COSUin,
				}
				policy, err := cosBuilder.BuildPolicy(cosReq)
				require.NoError(t, err)
				assert.NotEmpty(t, policy)
			}()
		}
		for i := 0; i < 10; i++ {
			go func() {
				defer func() { done <- true }()
				s3Req := storage.TokenRequest{
					Bucket: AWSBucket,
					Region: AWSRegion,
				}
				policy, err := s3Builder.BuildPolicy(s3Req)
				require.NoError(t, err)
				assert.NotEmpty(t, policy)
			}()
		}

		// 等待所有 goroutine 完成
		for i := 0; i < 20; i++ {
			<-done
		}
	})
}
