package provider

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/ratelimit"
)

type errMediaNotSupported struct {
	error
}

func mediaNotSupported(api string) errMediaNotSupported {
	return errMediaNotSupported{
		fmt.Errorf("media not supported %v", api),
	}
}

// IsErrMediaNotSupported return true if the error dictates taht the media
// type was not supported by the scraper
func IsErrMediaNotSupported(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(errMediaNotSupported); ok {
		return true
	}
	return false
}

// APIClient is a http client with rate limiting
type APIClient struct {
	*http.Client
	ratelimit.Limiter
}

// NewAPIClient return a new APIClient
func NewAPIClient(limit int) *APIClient {
	return &APIClient{
		Client:  http.DefaultClient,
		Limiter: ratelimit.New(limit),
	}
}

// Do performs a http request with rate limiting
func (a *APIClient) Do(req *http.Request) (*http.Response, error) {
	a.Take()
	return a.Client.Do(req)
}

// Get performs a http get request with rate limiting
func (a *APIClient) Get(url string) (*http.Response, error) {
	a.Take()
	return a.Client.Get(url)
}

// Head performs a http head request with rate limiting
func (a *APIClient) Head(url string) (*http.Response, error) {
	a.Take()
	return a.Client.Head(url)
}

// Post performs a http post request with rate limiting
func (a *APIClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	a.Take()
	return a.Client.Post(url, contentType, body)
}

// PostForm performs a http post form request with rate limiting
func (a *APIClient) PostForm(url string, data url.Values) (*http.Response, error) {
	a.Take()
	return a.Client.PostForm(url, data)
}
