package config

import (
	"github.com/defan6/http-checker/internal/checker"
	"github.com/defan6/http-checker/internal/reader"
	"github.com/defan6/http-checker/internal/writer"
)

const (
	outputFormates = 3
	inputFormates  = 2
)

func FormatInput(inFormat InputFormat) []func() []string {
	sliceFunc := make([]func() []string, 0, inputFormates)

	if inFormat.FileIn.IsFileIn {
		sliceFunc = append(sliceFunc, reader.FileReader(inFormat.FileIn.Path))
	}

	if inFormat.IsCmdLine || len(sliceFunc) == 0 {
		sliceFunc = append(sliceFunc, reader.CMDLineReader)
	}

	return sliceFunc
}

func FormatOutput(outFormat OutputFormat) []func([]checker.Result) {
	sliceFunc := make([]func([]checker.Result), 0, outputFormates)

	if outFormat.TXTFileOut.IsFileOut {
		sliceFunc = append(sliceFunc, writer.FormatFile(outFormat.TXTFileOut.Path))
	}

	if outFormat.IsJSONOut {
		sliceFunc = append(sliceFunc, writer.FormatJSON)
	}

	if outFormat.CSVFileOut.IsFileOut {
		sliceFunc = append(sliceFunc, writer.FormatCSVFile(outFormat.CSVFileOut.Path))
	}

	if outFormat.IsTableOut || len(sliceFunc) == 0 {
		sliceFunc = append(sliceFunc, writer.FormatTable)
	}

	return sliceFunc
}
