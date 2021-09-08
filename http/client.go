package http

import (
	"net/http"
)

type Client interface {
	Do(request *http.Request) (*http.Response, error)
}

type client struct {
	*http.Client
}

func NewClientWrapper(httpClient *http.Client) Client {
	return &client{
		Client: httpClient,
	}
}
