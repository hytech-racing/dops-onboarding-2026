package usecase

import (
	"context"
	"errors"
	"io"
	"main/internals/s3/repository"
)

// S3UseCaseInterface defines the business logic operations for S3 file management
type S3UseCaseInterface interface {
	// UploadFileFromReader uploads a file to S3 from an io.Reader
	UploadFileFromReader(ctx context.Context, reader io.Reader, objectName string) error
	// UploadFileFromWriter uploads a file to S3 from an io.WriterTo
	UploadFileFromWriter(ctx context.Context, writer *io.WriterTo, objectName string) error
	// GetPresignedUrl generates a presigned URL for accessing an S3 object
	GetPresignedUrl(ctx context.Context, objectPath string) (string, error)
}

// S3UseCase implements S3UseCaseInterface and handles S3 business logic
type S3UseCase struct {
	repo repository.S3RepositoryInterface
}

// NewS3UseCase creates a new S3UseCase instance with the given repository
func NewS3UseCase(repo repository.S3RepositoryInterface) *S3UseCase {
	return &S3UseCase{
		repo: repo,
	}
}

// UploadFileFromReader uploads a file to S3 from an io.Reader after validating inputs.
// Returns an error if the reader is nil or objectName is empty.
func (uc *S3UseCase) UploadFileFromReader(ctx context.Context, reader io.Reader, objectName string) error {
	if reader == nil {
		return errors.New("reader must not be nil")
	}

	if objectName == "" {
		return errors.New("objectName must not be empty")
	}

	return uc.repo.WriteObjectFromReader(ctx, reader, objectName)
}

// UploadFileFromWriter uploads a file to S3 from an io.WriterTo after validating inputs.
// Returns an error if the writer is nil or objectName is empty.
func (uc *S3UseCase) UploadFileFromWriter(ctx context.Context, writer *io.WriterTo, objectName string) error {
	if writer == nil {
		return errors.New("writer must not be nil")
	}

	if objectName == "" {
		return errors.New("objectName must not be empty")
	}

	return uc.repo.WriteObjectFromWriterTo(ctx, writer, objectName)
}

// GetPresignedUrl generates a presigned URL for accessing an S3 object after validating the objectPath.
// Returns an error if the objectPath is empty.
func (uc *S3UseCase) GetPresignedUrl(ctx context.Context, objectPath string) (string, error) {
	if objectPath == "" {
		return "", errors.New("objectPath must not be empty")
	}

	return uc.repo.GetSignedUrl(ctx, objectPath)
}
