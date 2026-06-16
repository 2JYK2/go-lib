package s3

import (
	"context"
	"errors"

	"github.com/2JYK2/go-lib/common/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

func (s *S3TokenManager) BatchDelete(ctx context.Context, keys []string) (*storage.DeleteResponse, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	if len(keys) > 1000 {
		return nil, errors.New("too many keys")
	}

	objects := make([]types.ObjectIdentifier, 0, len(keys))
	for _, k := range keys {
		objects = append(objects, types.ObjectIdentifier{
			Key: aws.String(k),
		})
	}

	resp, err := client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.BucketName),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return nil, err
	}

	if resp != nil {
		return &storage.DeleteResponse{
			Deleted: resp.Deleted,
			Errors:  resp.Errors,
		}, nil
	}

	return nil, nil
}
