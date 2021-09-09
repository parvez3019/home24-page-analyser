package http

import (
	"net/http"
)

type Client interface {
	Get(url string) (resp *http.Response, err error)
}

type client struct {
	*http.Client
}

func NewClientWrapper(httpClient *http.Client) Client {
	return &client{
		Client: httpClient,
	}
}
