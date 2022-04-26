package provider

import (
	"fmt"
	"github.com/klepatskiy/klepa-shop-backend/pkg/logger"
	"io"
	"net/http"
	"net/url"
)

type httpRequest struct {
	logger     logger.ZapLogger
	httpClient *http.Client
}

type Request struct {
	Method  string
	Url     string
	Body    io.Reader
	Params  map[string]string
	Headers map[string]string
	Cookies map[string]string
}

type HttpClient interface {
	MakeHttpRequest(r *Request) (*http.Response, error)
}

func NewHttpClient(zapLogger logger.ZapLogger) HttpClient {
	httpClient := &http.Client{}

	return &httpRequest{
		logger:     zapLogger,
		httpClient: httpClient,
	}
}

func (s *httpRequest) MakeHttpRequest(r *Request) (*http.Response, error) {
	u, _ := url.Parse(r.Url)
	urlParams := url.Values{}

	for key, value := range r.Params {
		urlParams.Add(key, value)
	}

	u.RawQuery = urlParams.Encode()
	requestUrl := u.String()

	if len(r.Method) == 0 {
		r.Method = http.MethodGet
	}

	req, err := http.NewRequest(r.Method, requestUrl, r.Body)
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	for key, value := range r.Cookies {
		cookie := http.Cookie{
			Name:  key,
			Value: value,
		}
		req.AddCookie(&cookie)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf(fmt.Sprintf("Error. Reponse code `%d` not equal `%d`", resp.StatusCode, http.StatusOK))
	}

	return resp, err
}
