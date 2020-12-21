package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"
)

type App interface {
	SQSWorker(ctx context.Context, message string) error
	S3Worker(ctx context.Context, bucket, filename string) error
	GetCat(ctx context.Context, id domain.CatID) (*domain.Cat, error)
	GetCats(ctx context.Context, first int64) (domain.Cats, error)
}

type Config struct {
	BucketName string
}

func NewApp(repos *domain.AllRepository, config *Config) App {
	return &app{
		repos:  repos,
		config: config,
	}
}

type app struct {
	repos  *domain.AllRepository
	config *Config
}

func (app *app) SQSWorker(ctx context.Context, message string) error {
	log.Debugf("log message %v", message)
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
		if err := app.repos.S3Client.Upload(ctx, app.config.BucketName, fmt.Sprintf("%v/%v.txt", message, cat.ID), file); err != nil {
			return err
		}
		log.Debugf("cat[%v] = %v", i, cat)
	}
	return nil
}

func (app *app) S3Worker(ctx context.Context, bucket, filename string) error {
	file, err := app.repos.S3Client.Download(ctx, bucket, filename)
	if err != nil {
		return nil
	}

	cat := &domain.Cat{}
	if err := json.Unmarshal(file, cat); err != nil {
		return err
	}

	if err := app.repos.Transaction(ctx, func(ctx context.Context, tx domain.Tx) error {
		if _, err := app.repos.CatRepository.CreateInTx(ctx, tx, cat); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (app *app) GetCat(ctx context.Context, id domain.CatID) (*domain.Cat, error) {
	cat, err := app.repos.CatRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (app *app) GetCats(ctx context.Context, first int64) (domain.Cats, error) {
	cats, err := app.repos.CatRepository.GetAll(ctx, first)
	if err != nil {
		return nil, err
	}
	return cats, nil
}
