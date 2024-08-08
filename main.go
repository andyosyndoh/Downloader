package main

import (
	"downloaderex/internal/downloader"
	"downloaderex/internal/fileManager"
	"downloaderex/internal/flags"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <URL> [options]")
		return
	}

	// Parse arguments
	args := os.Args[1:]
	url := ""
	file := ""
	rateLimit := ""

	for _, arg := range args {
		if strings.HasPrefix(arg, "--rate-limit=") {
			rateLimit = arg[len("--rate-limit="):]
		} else if strings.HasPrefix(arg, "-O=") {
			file = flags.GetFlagInput(arg)
		} else if strings.HasPrefix(arg, "-P=") {
			// Handle path if needed
			// You can combine path with file if necessary
		} else if strings.HasPrefix(arg, "http") {
			url = arg
		}
	}

	if url == "" {
		fmt.Println("Error: URL not provided.")
		return
	}

	// If no file is specified, derive it from the URL
	if file == "" {
		urlParts := strings.Split(url, "/")
		file = urlParts[len(urlParts)-1]
	}

	// Handle logger flag
	if len(os.Args) == 3 && os.Args[1] == "-B" {
		fileManager.Logger()
		fmt.Println("Output will be written to \"wget-log\"")
		return
	}

	// Start the download
	downloader.OneDownload(file, url, rateLimit)
}
