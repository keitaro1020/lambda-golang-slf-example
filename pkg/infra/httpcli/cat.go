package httpcli

import (
	"context"
	"net/http"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"
)

const (
	urlCatAPI = "https://api.thecatapi.com/v1/images/search"
)

type CatClient struct {
	httpClient HTTPClient
}

var _ domain.CatClient = &CatClient{}

func NewCatClient(httpClient HTTPClient) *CatClient {
	return &CatClient{
		httpClient: httpClient,
	}
}

func (cli *CatClient) Search(ctx context.Context) (domain.Cats, error) {
	cats := domain.Cats{}
	_, err := cli.httpClient.GetJSON(ctx, urlCatAPI, http.Header{}, nil, &cats)
	if err != nil {
		return nil, err
	}
	return cats, nil
}
