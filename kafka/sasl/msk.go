package sasl

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/aws"
)

func MskMechanism(creds *credentials.Credentials, userAgent string) sasl.Mechanism {
	return aws.ManagedStreamingIAM(func(ctx context.Context) (aws.Auth, error) {
		val, err := creds.GetWithContext(ctx)
		if err != nil {
			return aws.Auth{}, err
		}
		return aws.Auth{
			AccessKey:    val.AccessKeyID,
			SecretKey:    val.SecretAccessKey,
			SessionToken: val.SessionToken,
			UserAgent:    userAgent,
		}, nil
	})
}
