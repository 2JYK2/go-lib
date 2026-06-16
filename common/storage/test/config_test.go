package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfig 测试配置结构
type TestConfig struct {
	// AWS S3 配置
	AWS AWSConfig `json:"aws"`
	// COS 配置
	COS COSConfig `json:"cos"`
}

// AWSConfig AWS 配置
type AWSConfig struct {
	RoleArn         string `json:"roleArn"`
	Region          string `json:"region"`
	CallerSecretId  string `json:"callerSecretId"`
	CallerSecretKey string `json:"callerSecretKey"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
}

// COSConfig COS 配置
type COSConfig struct {
	CallerSecretId  string `json:"callerSecretId"`
	CallerSecretKey string `json:"callerSecretKey"`
	RoleArn         string `json:"roleArn"`
	Region          string `json:"region"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Uin             string `json:"uin"`
}

// GetTestConfig 获取测试配置
func GetTestConfig() *TestConfig {
	return &TestConfig{
		AWS: AWSConfig{
			RoleArn:         AWSRoleArn,
			Region:          AWSRegion,
			CallerSecretId:  AWSCallerSecretId,
			CallerSecretKey: AWSCallerSecretKey,
			Endpoint:        AWSEndpoint,
			Bucket:          AWSBucket,
		},
		COS: COSConfig{
			CallerSecretId:  COSCallerSecretId,
			CallerSecretKey: COSCallerSecretKey,
			RoleArn:         COSRoleArn,
			Region:          COSRegion,
			Endpoint:        COSEndpoint,
			Bucket:          COSBucket,
			Uin:             COSUin,
		},
	}
}

// TestConfigValidation 测试配置验证
func TestConfigValidation(t *testing.T) {
	config := GetTestConfig()

	t.Run("AWSConfigValidation", func(t *testing.T) {
		aws := config.AWS
		assert.NotEmpty(t, aws.RoleArn, "AWS RoleArn 不能为空")
		assert.NotEmpty(t, aws.Region, "AWS Region 不能为空")
		assert.NotEmpty(t, aws.CallerSecretId, "AWS CallerSecretId 不能为空")
		assert.NotEmpty(t, aws.CallerSecretKey, "AWS CallerSecretKey 不能为空")
		assert.NotEmpty(t, aws.Endpoint, "AWS Endpoint 不能为空")
		assert.NotEmpty(t, aws.Bucket, "AWS Bucket 不能为空")

		// 验证格式
		assert.Contains(t, aws.RoleArn, "arn:aws:iam::", "AWS RoleArn 格式不正确")
		assert.Contains(t, aws.CallerSecretId, "AKIA", "AWS CallerSecretId 格式不正确")
		assert.Len(t, aws.CallerSecretKey, 40, "AWS CallerSecretKey 长度不正确")
	})

	t.Run("COSConfigValidation", func(t *testing.T) {
		cos := config.COS
		assert.NotEmpty(t, cos.CallerSecretId, "COS CallerSecretId 不能为空")
		assert.NotEmpty(t, cos.CallerSecretKey, "COS CallerSecretKey 不能为空")
		assert.NotEmpty(t, cos.RoleArn, "COS RoleArn 不能为空")
		assert.NotEmpty(t, cos.Region, "COS Region 不能为空")
		assert.NotEmpty(t, cos.Endpoint, "COS Endpoint 不能为空")
		assert.NotEmpty(t, cos.Bucket, "COS Bucket 不能为空")
		assert.NotEmpty(t, cos.Uin, "COS Uin 不能为空")

		// 验证格式
		assert.Contains(t, cos.RoleArn, "qcs::cam::uin/", "COS RoleArn 格式不正确")
		assert.Contains(t, cos.CallerSecretId, "IKID", "COS CallerSecretId 格式不正确")
		assert.Len(t, cos.CallerSecretKey, 32, "COS CallerSecretKey 长度不正确")
		assert.Contains(t, cos.Endpoint, "cos.", "COS Endpoint 格式不正确")
	})

	t.Run("ConfigConsistency", func(t *testing.T) {
		// 验证配置一致性
		aws := config.AWS
		cos := config.COS

		// 验证区域格式
		assert.Regexp(t, `^[a-z]+-[a-z]+-[0-9]+$`, aws.Region, "AWS Region 格式不正确")
		assert.Regexp(t, `^[a-z]+-[a-z]+$`, cos.Region, "COS Region 格式不正确")

		// 验证存储桶名称格式
		assert.Regexp(t, `^[a-z0-9.-]+$`, aws.Bucket, "AWS Bucket 格式不正确")
		assert.Regexp(t, `^[a-z0-9.-]+$`, cos.Bucket, "COS Bucket 格式不正确")
	})
}

// TestEnvironmentConfig 测试环境变量配置
func TestEnvironmentConfig(t *testing.T) {
	t.Run("EnvironmentVariables", func(t *testing.T) {
		// 测试环境变量是否设置
		envVars := []string{
			"AWS_ACCESS_KEY_ID",
			"AWS_SECRET_ACCESS_KEY",
			"AWS_REGION",
			"AWS_S3_BUCKET",
		}

		for _, envVar := range envVars {
			value := os.Getenv(envVar)
			if value != "" {
				t.Logf("环境变量 %s 已设置", envVar)
			} else {
				t.Logf("环境变量 %s 未设置（这是正常的，因为使用硬编码配置）", envVar)
			}
		}
	})

	t.Run("ConfigOverride", func(t *testing.T) {
		// 测试配置覆盖
		config := GetTestConfig()

		// 验证默认配置
		assert.Equal(t, AWSRoleArn, config.AWS.RoleArn)
		assert.Equal(t, COSRoleArn, config.COS.RoleArn)

		// 这里可以添加环境变量覆盖的测试
		// 例如：如果设置了环境变量，应该使用环境变量的值
	})
}

// TestConfigSecurity 测试配置安全性
func TestConfigSecurity(t *testing.T) {
	t.Run("SecretMasking", func(t *testing.T) {
		config := GetTestConfig()

		// 验证密钥不为空
		assert.NotEmpty(t, config.AWS.CallerSecretKey)
		assert.NotEmpty(t, config.COS.CallerSecretKey)

		// 验证密钥长度
		assert.Len(t, config.AWS.CallerSecretKey, 40, "AWS 密钥长度应为 40")
		assert.Len(t, config.COS.CallerSecretKey, 32, "COS 密钥长度应为 32")

		// 注意：在实际测试中，不应该打印真实的密钥
		// 这里只是验证配置存在
		t.Log("配置安全性检查通过")
	})

	t.Run("ConfigValidation", func(t *testing.T) {
		config := GetTestConfig()

		// 验证配置完整性
		assert.NotNil(t, config)
		assert.NotNil(t, config.AWS)
		assert.NotNil(t, config.COS)

		// 验证必要字段
		requiredAWSFields := []string{
			config.AWS.RoleArn,
			config.AWS.Region,
			config.AWS.CallerSecretId,
			config.AWS.CallerSecretKey,
			config.AWS.Endpoint,
			config.AWS.Bucket,
		}

		requiredCOSFields := []string{
			config.COS.CallerSecretId,
			config.COS.CallerSecretKey,
			config.COS.RoleArn,
			config.COS.Region,
			config.COS.Endpoint,
			config.COS.Bucket,
			config.COS.Uin,
		}

		for i, field := range requiredAWSFields {
			assert.NotEmpty(t, field, "AWS 配置字段 %d 不能为空", i)
		}

		for i, field := range requiredCOSFields {
			assert.NotEmpty(t, field, "COS 配置字段 %d 不能为空", i)
		}
	})
}

// TestConfigCompatibility 测试配置兼容性
func TestConfigCompatibility(t *testing.T) {
	t.Run("CrossProviderCompatibility", func(t *testing.T) {
		config := GetTestConfig()

		// 验证两个提供商的配置都有效
		assert.NotEmpty(t, config.AWS.Bucket)
		assert.NotEmpty(t, config.COS.Bucket)
		assert.NotEqual(t, config.AWS.Bucket, config.COS.Bucket)

		assert.NotEmpty(t, config.AWS.Region)
		assert.NotEmpty(t, config.COS.Region)
		assert.NotEqual(t, config.AWS.Region, config.COS.Region)

		// 验证配置格式兼容性
		assert.Regexp(t, `^arn:aws:`, config.AWS.RoleArn)
		assert.Regexp(t, `^qcs::`, config.COS.RoleArn)
	})

	t.Run("ConfigConsistency", func(t *testing.T) {
		config := GetTestConfig()

		// 验证配置一致性
		assert.NotEmpty(t, config.AWS.CallerSecretId)
		assert.NotEmpty(t, config.COS.CallerSecretId)
		assert.NotEqual(t, config.AWS.CallerSecretId, config.COS.CallerSecretId)

		assert.NotEmpty(t, config.AWS.CallerSecretKey)
		assert.NotEmpty(t, config.COS.CallerSecretKey)
		assert.NotEqual(t, config.AWS.CallerSecretKey, config.COS.CallerSecretKey)
	})
}

// BenchmarkConfigAccess 配置访问性能测试
func BenchmarkConfigAccess(b *testing.B) {
	config := GetTestConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.AWS.RoleArn
		_ = config.AWS.Region
		_ = config.AWS.Bucket
		_ = config.COS.RoleArn
		_ = config.COS.Region
		_ = config.COS.Bucket
	}
}
