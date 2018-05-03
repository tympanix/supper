package provider

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/apex/log"
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
	name string
}

func simpleURL(url string) string {
	return strings.Split(url, "?")[0]
}

// NewAPIClient return a new APIClient
func NewAPIClient(name string, limit int) *APIClient {
	return &APIClient{
		Client:  http.DefaultClient,
		Limiter: ratelimit.New(limit),
		name:    name,
	}
}

// Do performs a http request with rate limiting
func (a *APIClient) Do(req *http.Request) (*http.Response, error) {
	a.Take()
	log.WithField("url", simpleURL(req.URL.String())).WithField("method", req.Method).Debug(a.name)
	return a.Client.Do(req)
}

// Get performs a http get request with rate limiting
func (a *APIClient) Get(url string) (*http.Response, error) {
	a.Take()
	log.WithField("url", simpleURL(url)).WithField("method", "GET").Debug(a.name)
	return a.Client.Get(url)
}

// Head performs a http head request with rate limiting
func (a *APIClient) Head(url string) (*http.Response, error) {
	a.Take()
	log.WithField("url", simpleURL(url)).WithField("method", "HEAD").Debug(a.name)
	return a.Client.Head(url)
}

// Post performs a http post request with rate limiting
func (a *APIClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	a.Take()
	log.WithField("url", simpleURL(url)).WithField("method", "POST").Debug(a.name)
	return a.Client.Post(url, contentType, body)
}

// PostForm performs a http post form request with rate limiting
func (a *APIClient) PostForm(url string, data url.Values) (*http.Response, error) {
	a.Take()
	log.WithField("url", simpleURL(url)).WithField("method", "POST").Debug(a.name)
	return a.Client.PostForm(url, data)
}
