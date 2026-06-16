package producer

import (
	"context"
	"crypto/tls"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/aws"
)

type Option interface {
	apply(config) config
}

type consumerOption func(config) config

func (o consumerOption) apply(conf config) config {
	return o(conf)
}

func DisableTls(e bool) Option {
	return consumerOption(func(c config) config {
		c.opts = append(c.opts, kgo.DialTLSConfig(&tls.Config{InsecureSkipVerify: e}))
		return c
	})
}

func WithLogger(logger kgo.Logger) Option {
	return consumerOption(func(c config) config {
		c.logger = logger
		return c
	})
}

func WithKgoOpts(opts ...kgo.Opt) Option {
	return consumerOption(func(cfg config) config {
		if len(opts) > 0 {
			cfg.opts = append(cfg.opts, opts...)
		}
		return cfg
	})
}

// WithSasl for all sasl
func WithSasl(mechanism sasl.Mechanism) Option {
	return consumerOption(func(cfg config) config {
		cfg.sasl = mechanism
		return cfg
	})
}

// WithAwsIAM only for aws iam auth
func WithAwsIAM(creds *credentials.Credentials, userAgent string) Option {
	return consumerOption(func(cfg config) config {
		cfg.sasl = aws.ManagedStreamingIAM(func(ctx context.Context) (aws.Auth, error) {
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
		return cfg
	})
}
