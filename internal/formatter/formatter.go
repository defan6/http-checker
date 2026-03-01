package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/defan6/http-checker/internal/checker"
)

type JSONResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Latency    string `json:"latency"`
	Error      string `json:"error,omitempty"` // добавляем поле для ошибки
}

func FormatTable(results []checker.Result) {
	if len(results) == 0 {
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Добавляем колонку Error
	fmt.Fprintln(w, "URL\tStatus Code\tLatency\tResult\tError")
	fmt.Fprintln(w, "---\t-----------\t-------\t------\t-----")

	for _, r := range results {
		statusSymbol := "❌"
		statusCode := r.StatusCode
		errorMsg := ""

		if r.Error != nil {
			statusSymbol = "💥"
			statusCode = 0
			errorMsg = r.Error.Error()
		} else if r.StatusCode >= 200 && r.StatusCode < 300 {
			statusSymbol = "✅"
		} else if r.StatusCode >= 300 && r.StatusCode < 400 {
			statusSymbol = "⚠️"
		}

		latencyStr := formatLatency(r.Latency)

		// Обрезаем слишком длинные ошибки для таблицы
		displayError := errorMsg
		if len(displayError) > 40 {
			displayError = displayError[:37] + "..."
		}

		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
			r.URL, statusCode, latencyStr, statusSymbol, displayError)
	}

	w.Flush()
}

func FormatJSON(results []checker.Result) {
	if len(results) == 0 {
		return
	}
	jsonResults := make([]JSONResult, 0, len(results))

	for _, res := range results {
		jsonResult := JSONResult{
			URL:        res.URL,
			StatusCode: res.StatusCode,
			Latency:    formatLatency(res.Latency),
		}

		// Добавляем ошибку в JSON если она есть
		if res.Error != nil {
			jsonResult.Error = res.Error.Error()
		}

		jsonResults = append(jsonResults, jsonResult)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(jsonResults); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding to JSON: %v\n", err)
	}
}

func FormatFile(path string) func([]checker.Result) {
	return func(results []checker.Result) {
		if len(results) == 0 {
			return
		}

		outputPath := path
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		if outputPath == "" {
			outputPath = filepath.Join("reports", fmt.Sprintf("check_%s.txt", timestamp))
		} else {
			// Исправляем формирование пути с таймштампом
			ext := filepath.Ext(outputPath)
			nameWithoutExt := strings.TrimSuffix(outputPath, ext)
			outputPath = filepath.Join("reports", fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext))
		}

		dir := filepath.Dir(outputPath)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				return
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			return
		}
		defer file.Close()

		w := tabwriter.NewWriter(file, 0, 0, 3, ' ', 0)

		fmt.Fprintln(w, "URL\tStatus Code\tLatency\tResult\tError")
		fmt.Fprintln(w, "---\t-----------\t-------\t------\t-----")

		for _, r := range results {
			statusSymbol := "❌"
			statusCode := r.StatusCode
			errorMsg := ""

			if r.Error != nil {
				statusSymbol = "💥"
				statusCode = 0
				errorMsg = r.Error.Error()
			} else if r.StatusCode >= 200 && r.StatusCode < 300 {
				statusSymbol = "✅"
			} else if r.StatusCode >= 300 && r.StatusCode < 400 {
				statusSymbol = "⚠️"
			}

			latencyStr := formatLatency(r.Latency)
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
				r.URL, statusCode, latencyStr, statusSymbol, errorMsg)
		}

		w.Flush()
		fmt.Printf("Results saved to: %s\n", outputPath)
	}
}

func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d µs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	} else {
		return fmt.Sprintf("%.2f s", d.Seconds())
	}
}
