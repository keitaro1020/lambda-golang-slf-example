package functions

import (
	"os"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/application"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/handler"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/infra"
)

func GetHandler() handler.Handler {
	return handler.NewHandler(
		application.NewApp(
			infra.NewAllRepository(),
			&application.Config{
				BucketName: os.Getenv("BucketName"),
			},
		),
	)
}
