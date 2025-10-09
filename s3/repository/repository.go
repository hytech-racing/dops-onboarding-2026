package repository

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3RepositoryInterface defines the contract for S3 operations
type S3RepositoryInterface interface {
	// WriteObjectFromWriterTo uploads an object to S3 from an io.WriterTo
	WriteObjectFromWriterTo(ctx context.Context, writer *io.WriterTo, objectName string) error
	// WriteObjectFromReader uploads an object to S3 from an io.Reader
	WriteObjectFromReader(ctx context.Context, reader io.Reader, objectName string) error
	// GetSignedUrl generates a presigned URL for accessing an S3 object
	GetSignedUrl(ctx context.Context, objectPath string) (string, error)
}

// S3Repository implements S3RepositoryInterface for AWS S3 operations
type S3Repository struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
}

// NewS3Repository creates a new S3Repository instance with the given S3 client and bucket name.
// Returns an error if the client is nil.
func NewS3Repository(client *s3.Client, bucketName string) (*S3Repository, error) {
	if client == nil {
		return nil, errors.New("client must not be nil")
	}
	return &S3Repository{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    bucketName,
	}, nil
}

// WriteObjectFromWriterTo uploads an object to S3 by first writing the content to a buffer from an io.WriterTo,
// then uploading the buffer to the specified object path in the bucket.
func (r *S3Repository) WriteObjectFromWriterTo(ctx context.Context, writer *io.WriterTo, objectName string) error {
	var buf bytes.Buffer

	_, err := (*writer).WriteTo(&buf)
	if err != nil {
		log.Printf("Failed to write buffer: %v", err)
	}

	reader := bytes.NewReader(buf.Bytes())
	_, err = r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(objectName),
		Body:   reader,
	})
	if err != nil {
		return errors.New("failed to upload file to s3")
	}

	return nil
}

// WriteObjectFromReader uploads an object to S3 directly from an io.Reader to the specified object path in the bucket.
func (r *S3Repository) WriteObjectFromReader(ctx context.Context, reader io.Reader, objectName string) error {
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(objectName),
		Body:   reader,
	})
	if err != nil {
		return errors.New("failed to upload file to s3")
	}

	return nil
}

// GetSignedUrl generates a presigned URL for accessing an S3 object.
// The URL expires after 10 minutes.
func (r *S3Repository) GetSignedUrl(ctx context.Context, objectPath string) (string, error) {
	request, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(objectPath),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(10 * int64(time.Minute))
	})
	if err != nil {
		return "", errors.New("failed to get presigned url")
	}

	return request.URL, nil
}
