package msk

import (
	"context"
	"crypto/tls"
	"github.com/IBM/sarama"
	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"io/ioutil"
)

type accessTokenProvider struct {
	awsRegion   string
	awsProfiles string
}

func (m *accessTokenProvider) Token() (*sarama.AccessToken, error) {
	if m.awsProfiles == "" {
		token, _, err := signer.GenerateAuthToken(context.Background(), m.awsRegion)
		return &sarama.AccessToken{Token: token}, err
	} else {
		token, _, err := signer.GenerateAuthTokenFromProfile(context.Background(), m.awsRegion, m.awsProfiles)
		return &sarama.AccessToken{Token: token}, err
	}
}

func InitConfig(awsRegion string, awsProfiles string, SASL, TLS bool) *sarama.Config {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	*/
	/*awsRegion, hasEnvRegion := os.LookupEnv(awsRegion)
	if !hasEnvRegion {
		log.Fatal("AWS_REGION environment variable not set")
	}*/

	configure := sarama.NewConfig()
	configure.Version = sarama.V3_5_1_0
	if SASL {
		configure.Net.SASL.Enable = true
		configure.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		configure.Net.SASL.TokenProvider = &accessTokenProvider{
			awsRegion:   awsRegion,
			awsProfiles: awsProfiles,
		}
	}

	if TLS {
		configure.Net.TLS.Enable = true
		configure.Net.TLS.Config = &tls.Config{}
	}
	configure.Consumer.Offsets.Initial = sarama.OffsetOldest
	return configure
}

func GetRockID() string {
	// Load default AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	// Create EC2 instance metadata service client
	imdsClient := imds.NewFromConfig(cfg)

	// Get availability zone ID
	output, err := imdsClient.GetMetadata(context.TODO(), &imds.GetMetadataInput{
		Path: "placement/availability-zone-id",
	})
	if err != nil {
		panic(err)
	}

	// Read content
	content, err := ioutil.ReadAll(output.Content)
	if err != nil {
		panic(err)
	}
	return string(content)
}
