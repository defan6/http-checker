package main

import (
	"flag"
	"fmt"

	"github.com/defan6/http-checker/internal/checker"
	"github.com/defan6/http-checker/internal/formatter"
	"github.com/defan6/http-checker/internal/runner"
)

const countFormates = 3

func main() {
	cfg := parseFlags()
	urls := flag.Args()
	if len(urls) > 0 {
		resultChan := make(chan checker.Result, len(urls))
		err := runner.Run(urls, cfg.Timeout, resultChan)
		if err != nil {
			fmt.Printf("Error processing url: %w", err)
			return
		}
		results := make([]checker.Result, 0, len(urls))
		for i := 0; i < len(urls); i++ {
			result := <-resultChan
			results = append(results, result)
		}

		outputWriters := formatOutput(cfg.OutFormat)
		write(results, outputWriters)
	}
}

type Config struct {
	Timeout   int
	OutFormat OutputFormat
	URLs      []string
}

type FileOut struct {
	IsFileOut bool
	Path      string
}

type OutputFormat struct {
	IsJSONOut  bool
	IsTableOut bool
	FileOut    FileOut
}

func formatOutput(outFormat OutputFormat) []func([]checker.Result) {
	sliceFunc := make([]func([]checker.Result), 0, countFormates)

	if outFormat.FileOut.IsFileOut {
		sliceFunc = append(sliceFunc, formatter.FormatFile(outFormat.FileOut.Path))
	}

	if outFormat.IsTableOut {
		sliceFunc = append(sliceFunc, formatter.FormatTable)
	}

	if outFormat.IsJSONOut {
		sliceFunc = append(sliceFunc, formatter.FormatJSON)
	}

	return sliceFunc
}

func write(res []checker.Result, writers []func([]checker.Result)) {
	for _, writer := range writers {
		writer(res)
	}
}

func parseFlags() Config {
	var cfg Config
	flag.IntVar(&cfg.Timeout, "t", 5, "request timeout in seconds")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "request timeout in seconds")
	flag.BoolVar(&cfg.OutFormat.IsJSONOut, "j", false, "json format output")
	flag.BoolVar(&cfg.OutFormat.IsJSONOut, "json", false, "json format output")
	flag.StringVar(&cfg.OutFormat.FileOut.Path, "o", "", "format file output")
	flag.StringVar(&cfg.OutFormat.FileOut.Path, "out", "", "format file output")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "o" {
			cfg.OutFormat.FileOut.IsFileOut = true
		}
	})
	fmt.Printf("Set timeout: %d\n", cfg.Timeout)
	fmt.Printf("Set file output: %v\n", cfg.OutFormat.FileOut.IsFileOut)
	fmt.Printf("Set table output: %v\n", cfg.OutFormat.IsTableOut)
	fmt.Printf("Set json output: %v\n", cfg.OutFormat.IsJSONOut)

	return cfg
}
