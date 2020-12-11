package infra

import (
	"os"

	"github.com/keitaro1020/lambda-golang-slf-example/service/domain"
	"github.com/keitaro1020/lambda-golang-slf-example/service/infra/db"
	"github.com/keitaro1020/lambda-golang-slf-example/service/infra/httpcli"
)

func NewAllRepository() *domain.AllRepository {
	httpCli := httpcli.NewHTTPClient()
	dbConfig := &db.Config{
		User:     os.Getenv("DB_USER"),
		Pass:     os.Getenv("DB_PASS"),
		Endpoint: os.Getenv("DB_ENDPOINT"),
		Name:     os.Getenv("DB_NAME"),
	}
	return &domain.AllRepository{
		CatClient:     httpcli.NewCatClient(httpCli),
		CatRepository: db.NewCatRepository(dbConfig),
		S3Client:      httpcli.NewS3Client(),
	}
}
