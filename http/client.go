package http

import (
	"net/http"
)

// Client Wrapper interface over http get method
type Client interface {
	Get(url string) (resp *http.Response, err error)
}

// client wrapper over http.Client
type client struct {
	*http.Client
}

// NewClientWrapper creates and return http client wrapper
func NewClientWrapper(httpClient *http.Client) Client {
	return &client{
		Client: httpClient,
	}
}
