package checker

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

func (c *Client) CheckURL(url string) Result {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	result := Result{URL: url}
	if err != nil {
		result.Error = err
		return result
	}
	start := time.Now()
	resp, err := c.httpClient.Do(req)
	requestTime := time.Since(start)
	result.Latency = requestTime
	if err != nil {
		result.Error = err
		return result
	}
	result.StatusCode = resp.StatusCode

	defer resp.Body.Close()

	return result
}
