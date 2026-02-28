package formatter

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/defan6/http-checker/internal/checker"
)

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

func formatLatency(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%d µs", d.Microseconds())
	} else if d < time.Second {
		return fmt.Sprintf("%d ms", d.Milliseconds())
	} else {
		return fmt.Sprintf("%.2f s", d.Seconds())
	}
}
