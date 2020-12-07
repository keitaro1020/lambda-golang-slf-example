package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/keitaro1020/lambda-golang-slf-example/service/domain"
)

type App interface {
	SQSWorker(ctx context.Context, message string) error
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
	log.Printf("log message %v", message)
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
		if err := app.repos.S3Client.Upload(ctx, os.Getenv("BucketName"), fmt.Sprintf("%v.txt", cat.ID), file); err != nil {
			return err
		}
		log.Printf("cat[%v] = %v", i, cat)
	}
	return nil
}
