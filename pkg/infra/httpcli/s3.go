package httpcli

import (
	"context"
	"io"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct{}

var _ domain.S3Client = &S3Client{}

func NewS3Client() *S3Client {
	return &S3Client{}
}

func (cli *S3Client) Upload(ctx context.Context, bucketName, key string, file io.Reader) error {
	uploader := s3manager.NewUploader(cli.newSession())
	if _, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	}); err != nil {
		return err
	}
	return nil
}

func (cli *S3Client) Download(ctx context.Context, bucketName, key string) ([]byte, error) {
	var buf []byte
	file := aws.NewWriteAtBuffer(buf)

	downloader := s3manager.NewDownloader(cli.newSession())
	if _, err := downloader.DownloadWithContext(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}); err != nil {
		return nil, err
	}
	return file.Bytes(), nil
}

func (cli *S3Client) newSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
	}))
}
