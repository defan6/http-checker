package client

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient http.Client
	Timeout    time.Duration
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		Timeout: timeout,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
