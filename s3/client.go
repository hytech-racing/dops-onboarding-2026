package s3

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client creates and returns a new AWS S3 client.
func NewS3Client(accessKey string, secretKey string, region string, endpoint string) (*s3.Client, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("accessKey and secretKey must not be empty")
	}

	if region == "" {
		return nil, errors.New("region must not be empty")
	}

	if endpoint == "" {
		return nil, errors.New("endpoint must not be empty")
	}

	staticCreds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(staticCreds))

	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return client, nil
}
