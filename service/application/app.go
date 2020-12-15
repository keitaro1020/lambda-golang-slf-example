package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/keitaro1020/lambda-golang-slf-example/service/domain"
)

type App interface {
	SQSWorker(ctx context.Context, message string) error
	S3Worker(ctx context.Context, bucket, filename string) error
}

func NewApp(repos *domain.AllRepository) App {
	return &app{
		repos: repos,
	}
}

type app struct {
	repos *domain.AllRepository
}

func (app *app) SQSWorker(ctx context.Context, message string) error {
	log.Debug("log message %v", message)
	cats, err := app.repos.CatClient.Search(ctx)
	if err != nil {
		return err
	}
	for i := range cats {
		cat := cats[i]
		data, err := json.Marshal(cat)
		if err != nil {
			return err
		}

		file := bytes.NewBuffer(data)
		if err := app.repos.S3Client.Upload(ctx, os.Getenv("BucketName"), fmt.Sprintf("%v/%v.txt", message, cat.ID), file); err != nil {
			return err
		}
		log.Debug("cat[%v] = %v", i, cat)
	}
	return nil
}

func (app *app) S3Worker(ctx context.Context, bucket, filename string) error {
	file, err := app.repos.S3Client.Download(ctx, bucket, filename)
	if err != nil {
		return nil
	}

	log.Debug("%v", string(file))
	return nil
}
