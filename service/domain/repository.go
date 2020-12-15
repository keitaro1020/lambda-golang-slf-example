package domain

import (
	"context"
)

type AllRepository struct {
	CatClient     CatClient
	S3Client      S3Client
	Transaction   func(ctx context.Context, txFunc func(ctx context.Context, tx Tx) error) (err error)
	CatRepository CatRepository
}

type Tx interface {
	Executor() interface{}
}
