package main

import (
	"fmt"
	"os"
	"strings"

	"downloaderex/internal/background"
	"downloaderex/internal/downloader"
	"downloaderex/internal/fileManager"
	"downloaderex/internal/flags"
	"downloaderex/internal/mirror"
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
	sourcefile := ""
	var workInBackground bool = false
	var log bool = false
	mirroring := false

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
		} else if strings.HasPrefix(arg, "-i=") {
			sourcefile = flags.GetFlagInput(arg)
		} else if strings.HasPrefix(arg, "--mirror") {
			mirroring = true
			break
		}
	}

	// if url == "" && sourcefile == "" {
	// 	fmt.Println("Error: URL not provided.")
	// 	return
	// }

	if mirroring {
		url, flagInput := mirror.GetMirrorUrl(args)
		mirror.DownloadPage(url, flagInput)
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
	if sourcefile != "" {
		downloader.DownloadMultipleFiles(sourcefile, file, rateLimit, path)
		return
	}
	if url == "" {
		fmt.Println("Error: URL not provided.")
		return
	}
	// Start the download
	downloader.OneDownload(file, url, rateLimit, path)
}
