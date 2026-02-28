package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/defan6/http-checker/internal/checker"
	"github.com/defan6/http-checker/internal/formatter"
)

func main() {
	var timeout int
	var jsonOutput bool
	flag.IntVar(&timeout, "t", 5, "request timeout in seconds")
	flag.IntVar(&timeout, "timeout", 5, "request timeout in seconds")
	flag.BoolVar(&jsonOutput, "j", false, "json format output")
	flag.BoolVar(&jsonOutput, "json", false, "json format output")
	flag.Parse()
	fmt.Printf("Set timeout: %d\n", timeout)
	fmt.Printf("Set json output: %v\n", jsonOutput)
	client := checker.NewClient(time.Second * 15)
	results := make([]checker.Result, 0, 10)
	urls := flag.Args()
	if len(urls) > 0 {
		for _, arg := range urls {
			res, err := client.CheckURL(arg)
			if err != nil {
				fmt.Printf("Error checking url for url %s: %w", arg, err)
				return
			}
			results = append(results, res)
		}

		outputWriter := formatOutput(jsonOutput)
		outputWriter(results)
	}
}

func formatOutput(json bool) func([]checker.Result) {
	if json {
		return formatter.FormatJSON
	} else {
		return formatter.FormatTable
	}
}
