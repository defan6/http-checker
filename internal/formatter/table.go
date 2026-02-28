package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/defan6/http-checker/internal/checker"
)

type JSONResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Latency    string `json:"latency_sec"`
}

func FormatTable(results []checker.Result) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "URL\tStatus Code\tLatency\tResult")
	fmt.Fprintln(w, "---\t-----------\t-------\t-------")

	for _, r := range results {
		statusSymbol := "❌" // красный крестик по умолчанию
		if r.StatusCode >= 200 && r.StatusCode < 300 {
			statusSymbol = "✅" // зеленая галочка для успеха
		} else if r.StatusCode >= 300 && r.StatusCode < 400 {
			statusSymbol = "⚠️" // предупреждение для редиректов
		}
		latencyStr := formatLatency(r.Latency)
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", r.URL, r.StatusCode, latencyStr, statusSymbol)
	}

	w.Flush()
}

func FormatJSON(results []checker.Result) {
	jsonResults := make([]JSONResult, len(results))
	for _, res := range results {
		jsonResult := JSONResult{
			URL:        res.URL,
			StatusCode: res.StatusCode,
			Latency:    formatLatency(res.Latency),
		}
		jsonResults = append(jsonResults, jsonResult)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")  // добавляем отступы (pretty format)
	encoder.SetEscapeHTML(false) // отключаем экранирование HTML

	if err := encoder.Encode(jsonResults); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding to JSON: %v\n", err)
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
