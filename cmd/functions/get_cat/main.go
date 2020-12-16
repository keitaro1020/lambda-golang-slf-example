package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keitaro1020/lambda-golang-slf-example/service/application"
	"github.com/keitaro1020/lambda-golang-slf-example/service/handler"
	"github.com/keitaro1020/lambda-golang-slf-example/service/infra"
	"github.com/keitaro1020/lambda-golang-slf-example/service/infra/logger"
)

func main() {
	logger.SetLogger()
	h := handler.NewHandler(application.NewApp(infra.NewAllRepository()))
	lambda.Start(h.GetCat)
}
