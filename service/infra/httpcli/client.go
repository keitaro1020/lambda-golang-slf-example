package httpcli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	neturl "net/url"
	"time"
)

// HTTPClient is a HTTP client interface
type HTTPClient interface {
	PostJSON(ctx context.Context, url string, header http.Header, req interface{}, res interface{}) (int, error)
	GetJSON(ctx context.Context, url string, header http.Header, param map[string]interface{}, res interface{}) (int, error)
}

type httpClient struct {
	client *http.Client
}

// NewHTTPClient is Create New HTTPClient
func NewHTTPClient() HTTPClient {
	client := &http.Client{Timeout: time.Duration(1) * time.Minute}

	return &httpClient{
		client: client,
	}
}

func (ci *httpClient) PostJSON(c context.Context, url string, header http.Header, req interface{}, res interface{}) (int, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(reqJSON))
	r.Header = header
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(c)

	// for debug
	dumpReq, _ := httputil.DumpRequestOut(r, true)
	log.Printf("dumpReq : %s", dumpReq)

	response, err := ci.client.Do(r)
	if err != nil {
		if response == nil {
			return http.StatusInternalServerError, err
		}
		return response.StatusCode, err
	}
	defer response.Body.Close()

	dumpRes, err := httputil.DumpResponse(response, true)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	log.Printf("dumpRes : %s", dumpRes)

	if res != nil {
		resBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if err = json.Unmarshal(resBody, &res); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return response.StatusCode, nil
}

func (ci *httpClient) GetJSON(c context.Context, url string, header http.Header, param map[string]interface{}, res interface{}) (int, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r, _ := http.NewRequest(http.MethodGet, u.String(), nil)
	r.Header = header
	r.Header.Set("Content-Type", "application/json")

	qparam := r.URL.Query()
	if param != nil {
		for k, v := range param {
			qparam.Add(k, fmt.Sprint(v))
		}
	}
	r.URL.RawQuery = qparam.Encode()
	r = r.WithContext(c)

	dumpReq, _ := httputil.DumpRequestOut(r, false)
	log.Printf("dumpReq : %s", dumpReq)

	resp, err := ci.client.Do(r)
	if err != nil {
		if resp == nil {
			return http.StatusInternalServerError, err
		}
		return resp.StatusCode, err
	}
	defer resp.Body.Close()

	dumpRes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	log.Printf("dumpRes : %s", dumpRes)

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := json.Unmarshal(resBody, &res); err != nil {
		return http.StatusInternalServerError, err
	}

	return resp.StatusCode, nil
}
