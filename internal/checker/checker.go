package checker

import (
	"net/http"
	"sync"
	"time"

	"github.com/defan6/http-checker/internal/client"
)

type Checker struct {
	client *client.Client
}

func NewChecker(client *client.Client) *Checker {
	return &Checker{
		client: client,
	}
}

func (c *Checker) CheckURL(url string) Result {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	result := Result{URL: url}
	if err != nil {
		result.Error = err
		return result
	}
	start := time.Now()
	resp, err := c.client.Do(req)
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

func (c *Checker) CheckURLsConcurrency(urls []string) <-chan Result {

	resultChan := make(chan Result, len(urls))
	wg := sync.WaitGroup{}

	go func() {
		for _, url := range urls {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result := c.CheckURL(url)
				resultChan <- result
			}()
		}
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}
