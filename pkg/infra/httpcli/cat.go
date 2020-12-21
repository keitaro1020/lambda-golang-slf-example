package httpcli

import (
	"context"
	"net/http"

	"github.com/keitaro1020/lambda-golang-slf-practice/pkg/domain"
)

const (
	urlCatAPI = "https://api.thecatapi.com/v1/images/search"
)

type catClient struct {
	httpClient HTTPClient
}

func NewCatClient(httpClient HTTPClient) domain.CatClient {
	return &catClient{
		httpClient: httpClient,
	}
}

func (cli *catClient) Search(ctx context.Context) (domain.Cats, error) {
	cats := domain.Cats{}
	_, err := cli.httpClient.GetJSON(ctx, urlCatAPI, http.Header{}, nil, &cats)
	if err != nil {
		return nil, err
	}
	return cats, nil
}
