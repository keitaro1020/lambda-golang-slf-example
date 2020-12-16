package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/keitaro1020/lambda-golang-slf-example/service/domain"

	"github.com/aws/aws-lambda-go/events"
	"github.com/keitaro1020/lambda-golang-slf-example/service/application"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	Ping(ctx context.Context) (Response, error)
	SQSWorker(ctx context.Context, sqsEvent events.SQSEvent) error
	S3Worker(ctx context.Context, s3Event events.S3Event) error
	GetCat(ctx context.Context, req Request) (Response, error)
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

type Request events.APIGatewayProxyRequest

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

func (h *handler) GetCat(ctx context.Context, req Request) (Response, error) {
	var resValue interface{}
	res := Response{StatusCode: http.StatusOK}
	id, ok := req.PathParameters["id"]
	if ok {
		cat, err := h.app.GetCat(ctx, domain.CatID(id))
		if err != nil {
			return h.errorResponse(http.StatusInternalServerError, err), nil
		}
		resValue = cat
	} else {
		cats, err := h.app.GetCats(ctx)
		if err != nil {
			return h.errorResponse(http.StatusInternalServerError, err), nil
		}
		resValue = cats
	}

	resBody, err := h.jsonString(resValue)
	if err != nil {
		return h.errorResponse(http.StatusInternalServerError, err), nil
	}
	res.Body = resBody

	return res, nil
}

func (h *handler) errorResponse(status int, err error) Response {
	body, err := h.jsonString(map[string]interface{}{
		"message": "Okay so your other function also executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: http.StatusInternalServerError, Body: err.Error()}
	}
	return Response{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func (h *handler) jsonString(v interface{}) (string, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)
	return buf.String(), nil
}
