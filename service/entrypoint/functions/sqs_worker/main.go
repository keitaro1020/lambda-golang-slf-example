package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keitaro1020/lambda-golang-slf-example/service/application"
	"github.com/keitaro1020/lambda-golang-slf-example/service/handler"
	"github.com/keitaro1020/lambda-golang-slf-example/service/infra"
)

func main() {
	h := handler.NewHandler(application.NewApp(infra.NewAllRepository()))
	lambda.Start(h.SQSWorker)
}
