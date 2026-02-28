package main

import (
	"fmt"
	"os"
	"time"

	"github.com/defan6/http-checker/internal/checker"
	"github.com/defan6/http-checker/internal/formatter"
)

func main() {
	client := checker.NewClient(time.Second * 5)
	results := make([]checker.Result, 0, 10)
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			res, err := client.CheckURL(arg)
			if err != nil {
				fmt.Printf("Error checking url for url %s: %w", arg, err)
				return
			}
			results = append(results, res)
		}

		formatter.FormatTable(results)
	}

}
