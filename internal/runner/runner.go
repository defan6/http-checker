package runner

import (
	"errors"
	"time"

	"github.com/defan6/http-checker/internal/checker"
	"github.com/defan6/http-checker/internal/client"
	"github.com/defan6/http-checker/internal/config"
	"github.com/defan6/http-checker/internal/reader"
	"github.com/defan6/http-checker/internal/writer"
)

func Run(cfg *config.Config) error {
	readers := config.FormatInput(cfg.InputFormat)
	urls := reader.ReadURLs(readers)
	if len(urls) == 0 {
		return errors.New("Zero urls recieved")
	}

	clTimeout := time.Duration(cfg.Timeout) * time.Second
	cl := client.NewClient(clTimeout)
	ch := checker.NewChecker(cl)

	resChan := ch.CheckURLsConcurrency(urls)

	results := make([]checker.Result, 0, len(urls))
	for v := range resChan {
		results = append(results, v)
	}

	writers := config.FormatOutput(cfg.OutFormat)
	writer.WriteURLs(results, writers)

	return nil
}
