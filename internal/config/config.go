package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Timeout     int
	MaxConnects int
	OutFormat   OutputFormat
	InputFormat InputFormat
	URLs        []string
}

type FileOut struct {
	IsFileOut bool
	Path      string
}

type InputFormat struct {
	IsCmdLine bool
	FileIn    FileIn
}

type FileIn struct {
	IsFileIn bool
	Path     string
}

type OutputFormat struct {
	IsJSONOut  bool
	IsTableOut bool
	CSVFileOut FileOut
	TXTFileOut FileOut
}

func ParseFlags() Config {
	var cfg Config
	flag.IntVar(&cfg.Timeout, "t", 5, "request timeout in seconds")
	flag.IntVar(&cfg.Timeout, "timeout", 5, "request timeout in seconds")
	flag.BoolVar(&cfg.OutFormat.IsJSONOut, "j", false, "json format output")
	flag.BoolVar(&cfg.OutFormat.IsJSONOut, "json", false, "json format output")
	flag.StringVar(&cfg.OutFormat.TXTFileOut.Path, "o", "", "format txt file output")
	flag.StringVar(&cfg.OutFormat.TXTFileOut.Path, "out", "", "format txt file output")
	flag.StringVar(&cfg.OutFormat.CSVFileOut.Path, "c", "", "format csv file output")
	flag.StringVar(&cfg.OutFormat.CSVFileOut.Path, "csv", "", "format csv file output")
	flag.IntVar(&cfg.MaxConnects, "max_conn", 10, "max connects")

	flag.StringVar(&cfg.InputFormat.FileIn.Path, "i", "", "input file with urls")
	flag.StringVar(&cfg.InputFormat.FileIn.Path, "in", "", "input file with urls")
	flag.BoolVar(&cfg.InputFormat.IsCmdLine, "u", false, "cmd line urls")
	flag.BoolVar(&cfg.InputFormat.IsCmdLine, "url", false, "cmd line urls")
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "o", "out":
			cfg.OutFormat.TXTFileOut.IsFileOut = true
		case "c", "csv":
			cfg.OutFormat.CSVFileOut.IsFileOut = true
		case "i", "in":
			cfg.InputFormat.FileIn.IsFileIn = true
		}
	})

	Print(cfg)

	return cfg
}

func Print(cfg Config) {
	fmt.Printf("Set timeout: %d\n", cfg.Timeout)
	fmt.Printf("Set txt file input: %v\n", cfg.InputFormat.FileIn.IsFileIn)
	fmt.Printf("Set txt file output: %v\n", cfg.OutFormat.TXTFileOut.IsFileOut)
	fmt.Printf("Set csv file output: %v\n", cfg.OutFormat.CSVFileOut.IsFileOut)
	fmt.Printf("Set table output: %v\n", cfg.OutFormat.IsTableOut)
	fmt.Printf("Set json output: %v\n", cfg.OutFormat.IsJSONOut)
	fmt.Printf("Set input file: %v\n", cfg.InputFormat.FileIn.IsFileIn)
	fmt.Printf("Set cmd line input: %v\n", cfg.InputFormat.IsCmdLine)
	fmt.Printf("Set max connects: %d\n", cfg.MaxConnects)
}
