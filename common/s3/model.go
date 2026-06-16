package s3

type UploadConfig struct {
	AccessKeyId     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`
	RegionName      string `json:"regionName"`
	BucketName      string `json:"bucketName"`
}
type UploadRequest struct {
	ServerName string `json:"serverName"`
	PolicyUrl  string `json:"policyUrl"`
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	Duration   int64  `json:"duration"`
	RoleArn    string `json:"roleArn"`
}

type AWSPolicy struct {
	Version   string               `json:"Version"`
	Statement []AWSPolicyStatement `json:"Statement"`
}
type AWSPolicyStatement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource string   `json:"Resource"`
}
