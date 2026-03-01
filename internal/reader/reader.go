package reader

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func FileReader(path string) func() []string {
	return func() []string {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "path %s does not exists: %v\n", path, err)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening file by path %s: %v\n", path, err)
			return nil
		}

		defer file.Close()

		var urls []string
		scanner := bufio.NewScanner(file)
		lineNum := 0

		for scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())

			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			urls = append(urls, line)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading file at line %d: %w", lineNum, err)
			return nil
		}

		return urls
	}
}

func CMDLineReader() []string {
	fmt.Printf("Flag args: %s", flag.Args())
	return flag.Args()
}

func ReadURLs(readers []func() []string) []string {
	urls := make([]string, 0, 100)
	for _, reader := range readers {
		if reader == nil {
			continue
		}
		urls = append(urls, reader()...)
	}

	return urls
}
