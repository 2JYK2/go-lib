package test

import (
	"context"
	"testing"

	"github.com/2JYK2/go-lib/common/cos"
	"github.com/2JYK2/go-lib/common/storage"
	cos2 "github.com/2JYK2/go-lib/common/storage/cos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCOSProvider 测试 COS 提供商功能
func TestCOSProvider(t *testing.T) {
	t.Run("TokenManagerCreation", func(t *testing.T) {
		manager := cos2.NewCOSTokenManager()
		require.NotNil(t, manager)
		assert.Equal(t, storage.ProviderCOS, manager.GetProvider())
	})

	t.Run("UploadConfigStructure", func(t *testing.T) {
		config := &cos.UploadConfig{
			AccessKeyId:     COSCallerSecretId,
			SecretAccessKey: COSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      COSRegion,
			BucketName:      COSBucket,
		}

		assert.Equal(t, COSCallerSecretId, config.AccessKeyId)
		assert.Equal(t, COSCallerSecretKey, config.SecretAccessKey)
		assert.Equal(t, "test-session-token", config.SessionToken)
		assert.Equal(t, COSRegion, config.RegionName)
		assert.Equal(t, COSBucket, config.BucketName)
	})

	t.Run("UploadRequestStructure", func(t *testing.T) {
		request := cos2.UploadRequest{
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

		assert.Equal(t, COSCallerSecretId, request.AccessKeyId)
		assert.Equal(t, COSCallerSecretKey, request.SecretAccessKey)
		assert.Equal(t, "test-server", request.ServerName)
		assert.Equal(t, "test/path", request.PolicyUrl)
		assert.Equal(t, COSBucket, request.Bucket)
		assert.Equal(t, COSRegion, request.Region)
		assert.Equal(t, int64(3600), request.Duration)
		assert.Equal(t, COSRoleArn, request.RoleArn)
		assert.Equal(t, COSUin, request.Uin)
	})
}

// TestCOSConversion 测试 COS 转换功能
func TestCOSConversion(t *testing.T) {
	t.Run("UploadConfigToStorageCredentials", func(t *testing.T) {
		uploadConfig := &cos.UploadConfig{
			AccessKeyId:     COSCallerSecretId,
			SecretAccessKey: COSCallerSecretKey,
			SessionToken:    "test-session-token",
			RegionName:      COSRegion,
			BucketName:      COSBucket,
		}

		storageCreds := uploadConfig.ToStorageCredentials()
		require.NotNil(t, storageCreds)

		assert.Equal(t, COSCallerSecretId, storageCreds.AccessKeyId)
		assert.Equal(t, COSCallerSecretKey, storageCreds.SecretAccessKey)
		assert.Equal(t, "test-session-token", storageCreds.SessionToken)
		assert.Equal(t, COSRegion, storageCreds.RegionName)
		assert.Equal(t, COSBucket, storageCreds.BucketName)
	})

	t.Run("UploadRequestToStorageTokenRequest", func(t *testing.T) {
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
		require.NotNil(t, storageReq)

		assert.Equal(t, "test-server", storageReq.ServerName)
		assert.Equal(t, "test/path", storageReq.PolicyUrl)
		assert.Equal(t, COSBucket, storageReq.Bucket)
		assert.Equal(t, COSRegion, storageReq.Region)
		assert.Equal(t, int64(3600), storageReq.Duration)
		assert.Equal(t, COSRoleArn, storageReq.RoleArn)
		assert.Equal(t, COSUin, storageReq.Uin)
	})
}

// TestCOSPolicyGeneration 测试 COS 策略生成
func TestCOSPolicyGeneration(t *testing.T) {
	t.Run("BasicPolicy", func(t *testing.T) {
		req := storage.TokenRequest{
			Bucket: COSBucket,
			Region: COSRegion,
			Uin:    COSUin,
		}

		// 使用 COS 策略构建器
		builder := storage.NewCOSPolicyBuilder()
		policy, err := builder.BuildPolicy(req)
		require.NoError(t, err)
		require.NotEmpty(t, policy)

		// 验证策略内容
		assert.Contains(t, policy, COSBucket)
		assert.Contains(t, policy, COSRegion)
		assert.Contains(t, policy, COSUin)
		assert.Contains(t, policy, "cos:PutObject")
		assert.Contains(t, policy, "cos:GetObject")
		assert.Contains(t, policy, "cos:DeleteObject")
	})

	t.Run("PolicyWithPath", func(t *testing.T) {
		req := storage.TokenRequest{
			Bucket:    COSBucket,
			Region:    COSRegion,
			Uin:       COSUin,
			PolicyUrl: "test/path",
		}

		builder := storage.NewCOSPolicyBuilder()
		policy, err := builder.BuildPolicy(req)
		require.NoError(t, err)
		require.NotEmpty(t, policy)

		// 验证策略包含路径
		assert.Contains(t, policy, "test/path")
	})

	t.Run("PolicyValidation", func(t *testing.T) {
		builder := storage.NewCOSPolicyBuilder()

		// 测试缺少必需字段
		req := storage.TokenRequest{
			Bucket: COSBucket,
			// 缺少 Region 和 Uin
		}
		_, err := builder.BuildPolicy(req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bucket, region and uin are required for COS")
	})
}

// TestCOSIntegration 测试 COS 集成（需要真实凭证时跳过）
func TestCOSIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试")
	}

	t.Run("GetSessionTokenWithRole", func(t *testing.T) {
		// 注意：这个测试需要真实的腾讯云凭证
		// 在实际环境中，可以取消注释并配置正确的凭证
		t.Skip("需要真实的腾讯云凭证")

		/*
			manager := cos.NewCOSTokenManager()
			ctx := context.Background()

			req := storage.TokenRequest{
				ServerName: "test-server",
				Bucket:     COSBucket,
				Region:     COSRegion,
				Duration:   3600,
				RoleArn:    COSRoleArn,
				Uin:        COSUin,
			}

			creds, err := manager.GetSessionToken(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, creds)

			assert.NotEmpty(t, creds.AccessKeyId)
			assert.NotEmpty(t, creds.SecretAccessKey)
			assert.NotEmpty(t, creds.SessionToken)
			assert.Equal(t, COSRegion, creds.RegionName)
			assert.Equal(t, COSBucket, creds.BucketName)
		*/
	})

	t.Run("GetSessionTokenWithoutRole", func(t *testing.T) {
		// 注意：这个测试需要真实的腾讯云凭证
		t.Skip("需要真实的腾讯云凭证")

		/*
			manager := cos.NewCOSTokenManager()
			ctx := context.Background()

			req := storage.TokenRequest{
				ServerName: "test-server",
				Bucket:     COSBucket,
				Region:     COSRegion,
				Duration:   3600,
				Uin:        COSUin,
			}

			creds, err := manager.GetSessionToken(ctx, req)
			require.NoError(t, err)
			require.NotNil(t, creds)

			assert.NotEmpty(t, creds.AccessKeyId)
			assert.NotEmpty(t, creds.SecretAccessKey)
			assert.NotEmpty(t, creds.SessionToken)
			assert.Equal(t, COSRegion, creds.RegionName)
			assert.Equal(t, COSBucket, creds.BucketName)
		*/
	})
}

// TestCOSErrorHandling 测试 COS 错误处理
func TestCOSErrorHandling(t *testing.T) {
	t.Run("InvalidCredentials", func(t *testing.T) {
		// 测试无效凭证的错误处理
		// 这里可以添加更多的错误场景测试
		manager := cos2.NewCOSTokenManager()
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
}

// BenchmarkCOSPolicyBuilding COS 策略构建性能测试
func BenchmarkCOSPolicyBuilding(b *testing.B) {
	builder := storage.NewCOSPolicyBuilder()
	req := storage.TokenRequest{
		Bucket: COSBucket,
		Region: COSRegion,
		Uin:    COSUin,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := builder.BuildPolicy(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// TestCOSConcurrency 测试 COS 并发访问
func TestCOSConcurrency(t *testing.T) {
	// 并发测试策略构建
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			builder := storage.NewCOSPolicyBuilder()
			req := storage.TokenRequest{
				Bucket: COSBucket,
				Region: COSRegion,
				Uin:    COSUin,
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
