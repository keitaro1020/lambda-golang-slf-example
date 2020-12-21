package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/application"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/handler"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/infra"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/infra/logger"
)

func main() {
	logger.SetLogger()
	h := handler.NewHandler(application.NewApp(infra.NewAllRepository()))
	lambda.Start(h.S3Worker)
}
