package test

import (
	"os"
)

// 测试配置常量 - 从环境变量获取
var (
	// AWS S3 配置
	AWSCallerSecretId  = os.Getenv("AWS_CALLER_SECRET_ID")
	AWSCallerSecretKey = os.Getenv("AWS_CALLER_SECRET_KEY")
	AWSRoleArn         = os.Getenv("AWS_ROLE_ARN")
	AWSRegion          = os.Getenv("AWS_REGION")
	AWSEndpoint        = os.Getenv("AWS_ENDPOINT")
	AWSBucket          = os.Getenv("AWS_BUCKET")

	// COS 配置
	COSCallerSecretId  = os.Getenv("COS_CALLER_SECRET_ID")
	COSCallerSecretKey = os.Getenv("COS_CALLER_SECRET_KEY")
	COSRoleArn         = os.Getenv("COS_ROLE_ARN")
	COSRegion          = os.Getenv("COS_REGION")
	COSEndpoint        = os.Getenv("COS_ENDPOINT")
	COSBucket          = os.Getenv("COS_BUCKET")
	COSUin             = os.Getenv("COS_UIN")
)
