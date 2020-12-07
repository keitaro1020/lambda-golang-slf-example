package httpcli

import (
	"context"
	"io"

	"github.com/keitaro1020/lambda-golang-slf-example/service/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws/session"
)

type s3Client struct{}

func NewS3Client() domain.S3Client {
	return &s3Client{}
}

func (cli *s3Client) Upload(ctx context.Context, bucketName, key string, file io.Reader) error {
	uploader := s3manager.NewUploader(cli.newSession())
	if _, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	}); err != nil {
		return err
	}
	return nil
}

func (cli *s3Client) newSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
}
