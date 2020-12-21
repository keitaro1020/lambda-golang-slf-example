package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keitaro1020/lambda-golang-slf-practice/service/application"
	"github.com/keitaro1020/lambda-golang-slf-practice/service/handler"
	"github.com/keitaro1020/lambda-golang-slf-practice/service/infra"
	"github.com/keitaro1020/lambda-golang-slf-practice/service/infra/logger"
)

func main() {
	logger.SetLogger()
	h := handler.NewHandler(application.NewApp(infra.NewAllRepository()))
	lambda.Start(h.Graphql)
}
