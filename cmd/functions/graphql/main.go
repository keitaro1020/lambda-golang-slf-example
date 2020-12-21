package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/keitaro1020/lambda-golang-slf-practice/cmd/functions"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/infra/logger"
)

func main() {
	logger.SetLogger()
	lambda.Start(functions.GetHandler().Graphql)
}
