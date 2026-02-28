package checker

import (
	"fmt"
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

func (c *Client) CheckURL(url string) (Result, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Result{}, fmt.Errorf("Cannot create req for url %s: %w", url, err)
	}
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	requestTime := time.Since(start)
	if err != nil {
		return Result{}, fmt.Errorf("Cannot send req for url %s: %w", url, err)
	}

	defer resp.Body.Close()
	res := Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Latency:    requestTime,
	}

	return res, nil
}
