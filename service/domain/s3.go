package domain

import (
	"context"
	"io"
)

type S3Client interface {
	Upload(ctx context.Context, bucketName, key string, file io.Reader) error
	Download(ctx context.Context, bucketName, key string) ([]byte, error)
}
