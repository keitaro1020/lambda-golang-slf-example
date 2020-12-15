package handler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/keitaro1020/lambda-golang-slf-example/service/application"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	Ping(ctx context.Context) (Response, error)
	SQSWorker(ctx context.Context, sqsEvent events.SQSEvent) error
	S3Worker(ctx context.Context, s3Event events.S3Event) error
}

func NewHandler(app application.App) Handler {
	return &handler{
		app: app,
	}
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type handler struct {
	app application.App
}

func (h *handler) Ping(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Okay so your other function also executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "ping-cmd",
		},
	}
	log.Infof("Ping Response: %v", &resp)

	return resp, nil
}

func (h *handler) SQSWorker(ctx context.Context, sqsEvent events.SQSEvent) error {
	// todo 入力パラメータのチェック
	for i := range sqsEvent.Records {
		if err := h.app.SQSWorker(ctx, sqsEvent.Records[i].Body); err != nil {
			log.Errorf("SQSWorker error = %v", err)
			return err
		}
	}
	return nil
}

func (h *handler) S3Worker(ctx context.Context, s3Event events.S3Event) error {
	// todo 入力パラメータのチェック
	log.Debugf("%+v", s3Event)
	for i := range s3Event.Records {
		s3 := s3Event.Records[i].S3
		if err := h.app.S3Worker(ctx, s3.Bucket.Name, s3.Object.Key); err != nil {
			log.Errorf("SQSWorker error = %v", err)
			return err
		}
	}
	return nil
}
