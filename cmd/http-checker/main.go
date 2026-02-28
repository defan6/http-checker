package main

import (
	"fmt"
	"os"
	"time"

	"github.com/defan6/http-checker/internal/checker"
)

func main() {
	client := checker.NewClient(time.Second * 10)
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			res, err := client.CheckURL(arg)
			if err != nil {
				fmt.Printf("Error checking url for url %s: %w", arg, res)
				return
			}
			fmt.Printf("Response from %s: %+v\n", arg, res)
		}
	}
}
