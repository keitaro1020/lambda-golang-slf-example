package handler

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/keitaro1020/lambda-golang-slf-example/service/application"

	"github.com/aws/aws-lambda-go/events"
)

type Handler interface {
	Ping(ctx context.Context) (Response, error)
	SQSWorker(ctx context.Context, sqsEvent events.SQSEvent) error
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

	return resp, nil
}

func (h *handler) SQSWorker(ctx context.Context, sqsEvent events.SQSEvent) error {
	for i := range sqsEvent.Records {
		if err := h.app.SQSWorker(ctx, sqsEvent.Records[i].Body); err != nil {
			return err
		}
	}
	return nil
}
