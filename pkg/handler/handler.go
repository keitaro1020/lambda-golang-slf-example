package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	gqlgenhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aws/aws-lambda-go/events"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/application"
	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"
	"github.com/keitaro1020/lambda-golang-slf-practice/scripts/graphql/generated"
)

type Handler interface {
	Ping(ctx context.Context) (events.APIGatewayProxyResponse, error)
	SQSWorker(ctx context.Context, sqsEvent events.SQSEvent) error
	S3Worker(ctx context.Context, s3Event events.S3Event) error
	GetCat(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	Graphql(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewHandler(app application.App) Handler {
	return &handler{
		app: app,
	}
}

type handler struct {
	app application.App
}

func (h *handler) Ping(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Okay so your other function also executed successfully!",
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := events.APIGatewayProxyResponse{
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

func (h *handler) GetCat(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var resValue interface{}
	res := events.APIGatewayProxyResponse{StatusCode: http.StatusOK}
	id, ok := req.PathParameters["id"]
	if ok {
		cat, err := h.app.GetCat(ctx, domain.CatID(id))
		if err != nil {
			return h.errorResponse(http.StatusInternalServerError, err), nil
		}
		resValue = cat
	} else {
		cats, err := h.app.GetCats(ctx, 0)
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

var chiLambda *chiadapter.ChiLambda

func (h *handler) Graphql(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Infof("request: %#v", req)
	if chiLambda == nil {
		r := chi.NewRouter()
		r.Route("/graphql", func(r chi.Router) {
			r.Post("/query", gqlgenhandler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
				Resolvers: NewResolver(h.app),
			})).ServeHTTP)
			r.Get("/playground", playground.Handler("GraphQL", "/dev/graphql/query"))
			r.Get("/ping", func(w http.ResponseWriter, req *http.Request) {
				payload := struct {
					Message string
				}{
					Message: "pong",
				}
				res, err := json.Marshal(payload)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				w.Write(res)
			})
		})
		chiLambda = chiadapter.New(r)
	}

	return chiLambda.ProxyWithContext(ctx, req)
}

func (h *handler) errorResponse(status int, err error) events.APIGatewayProxyResponse {
	body, err := h.jsonString(map[string]interface{}{
		"message": "Okay so your other function also executed successfully!",
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: err.Error()}
	}
	return events.APIGatewayProxyResponse{
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
