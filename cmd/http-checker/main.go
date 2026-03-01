package main

import (
	"fmt"
	"os"

	"github.com/defan6/http-checker/internal/config"
	"github.com/defan6/http-checker/internal/runner"
)

func main() {
	cfg := config.ParseFlags()
	if err := runner.Run(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error(main): %w", err)
	}
}
