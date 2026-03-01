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

func (c *Checker) CheckURLsConcurrency(maxConnects int, urls []string) <-chan Result {
	urlsCh := generate(urls)
	resultChan := make(chan Result, len(urls))
	wg := sync.WaitGroup{}

	for range maxConnects {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range urlsCh {
				result := c.CheckURL(url)
				resultChan <- result
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

func generate(urls []string) chan string {
	in := make(chan string)

	go func() {
		for _, url := range urls {
			in <- url
		}
		close(in)
	}()

	return in
}
