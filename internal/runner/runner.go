package runner

import (
	"sync"
	"time"

	"github.com/defan6/http-checker/internal/checker"
)

func Run(urls []string, timeout int, resCh chan checker.Result) error {
	clTimeout := time.Duration(timeout) * time.Second
	client := checker.NewClient(clTimeout)
	wg := sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := client.CheckURL(url)
			resCh <- result
		}()
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	return nil
}
