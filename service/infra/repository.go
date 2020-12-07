package infra

import (
	"github.com/keitaro1020/lambda-golang-slf-example/service/domain"
	"github.com/keitaro1020/lambda-golang-slf-example/service/infra/httpcli"
)

func NewAllRepository() *domain.AllRepository {
	httpCli := httpcli.NewHTTPClient()
	return &domain.AllRepository{
		CatClient: httpcli.NewCatClient(httpCli),
		S3Client:  httpcli.NewS3Client(),
	}
}
