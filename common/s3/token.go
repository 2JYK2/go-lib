package s3

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type CustomLogger struct{}

func (c *CustomLogger) Log(args ...interface{}) {
	log.Println(args...) // 将日志输出到标准输出
}

func GetSessionToken(ctx context.Context, uploadRequest UploadRequest) (*UploadConfig, error) {
	resourceBucket := fmt.Sprintf("arn:aws:s3:::%s", uploadRequest.Bucket)
	resource := fmt.Sprintf("arn:aws:s3:::%s/*", uploadRequest.Bucket)
	if uploadRequest.PolicyUrl != "" {
		resource = fmt.Sprintf("arn:aws:s3:::%s/%s/*", uploadRequest.Bucket, uploadRequest.PolicyUrl)
	}
	policy := AWSPolicy{
		Version: "2012-10-17",
		Statement: []AWSPolicyStatement{
			{
				Effect:   "Allow",
				Action:   []string{"s3:GetObject", "s3:CreateMultipartUpload", "s3:CompleteMultipartUpload", "s3:UploadPart", "s3:PutObject", "s3:PutObjectAcl"},
				Resource: resource,
			},
			{
				Effect:   "Allow",
				Action:   []string{"s3:List*"},
				Resource: resourceBucket,
			},
		},
	}
	policyBytes, _ := json.Marshal(policy)
	policyStr := string(policyBytes)

	// create session
	s1, err := session.NewSession(&aws.Config{
		Region: &uploadRequest.Region,
	})

	if err != nil {
		return nil, err
	}

	var uploadConfig *UploadConfig
	svc := sts.New(s1)
	if uploadRequest.RoleArn != "" {
		input := &sts.AssumeRoleInput{
			RoleArn:         aws.String(uploadRequest.RoleArn),
			RoleSessionName: aws.String(uploadRequest.ServerName),
			Policy:          aws.String(policyStr),
			DurationSeconds: aws.Int64(uploadRequest.Duration),
			Tags:            []*sts.Tag{},
		}

		result, err := svc.AssumeRole(input)
		if err != nil {
			return nil, err
		}
		uploadConfig = &UploadConfig{
			AccessKeyId:     *result.Credentials.AccessKeyId,
			SecretAccessKey: *result.Credentials.SecretAccessKey,
			SessionToken:    *result.Credentials.SessionToken,
			BucketName:      uploadRequest.Bucket,
			RegionName:      uploadRequest.Region,
		}
	} else {
		input := &sts.GetFederationTokenInput{
			DurationSeconds: aws.Int64(uploadRequest.Duration),
			Name:            aws.String(uploadRequest.ServerName),
			Policy:          aws.String(policyStr),
			Tags:            []*sts.Tag{},
		}
		result, err := svc.GetFederationToken(input)
		if err != nil {
			return nil, err
		}
		uploadConfig = &UploadConfig{
			AccessKeyId:     *result.Credentials.AccessKeyId,
			SecretAccessKey: *result.Credentials.SecretAccessKey,
			SessionToken:    *result.Credentials.SessionToken,
			BucketName:      uploadRequest.Bucket,
			RegionName:      uploadRequest.Region,
		}
	}

	return uploadConfig, nil
}
