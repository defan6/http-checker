package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/defan6/http-checker/internal/config"
	"github.com/defan6/http-checker/internal/runner"
)

func main() {
	cfg := config.ParseFlags()
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := runner.Run(&cfg)
		if err != nil {
			fmt.Printf("Error in main: %w", err)
		}
	}()

	<-sigchan
	fmt.Println("An interrupt signal was received. Terminating work...")
	time.Sleep(2 * time.Second)
	fmt.Println("Work finished successfully!")
}
