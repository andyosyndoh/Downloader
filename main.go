package main

import (
	"downloaderex/internal/background"
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
	path := ""
	var workInBackground bool = false
	var log bool = false

	for _, arg := range args {
		if strings.HasPrefix(arg, "--rate-limit=") {
			rateLimit = arg[len("--rate-limit="):]
		} else if strings.HasPrefix(arg, "-O=") {
			file = flags.GetFlagInput(arg)
		} else if strings.HasPrefix(arg, "-P=") {
			path = flags.GetFlagInput(arg)
		} else if strings.HasPrefix(arg, "http") {
			url = arg
		} else if strings.HasPrefix(arg, "-B") {
			workInBackground = true
		} else if arg == "-b" {
			log = true
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
	if workInBackground {
		background.DownloadInBackground(file, url, rateLimit)
		return
	}
	if log {
		fileManager.Logger(file, url, rateLimit)
		return
	}

	// Start the download
	downloader.OneDownload(file, url, rateLimit, path)
}
