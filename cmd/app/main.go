package main

import (
	"fmt"
	"os"
	"strings"

	"downloaderex/internal/background"
	"downloaderex/internal/downloader"
	"downloaderex/internal/flags"
	"downloaderex/internal/mirror"
)

func main() {
	// Check if arguments are provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <URL> [options]")
		return
	}

	inputs := flags.ParseArgs()

	// Mirror website handling
	if inputs.Mirroring {
		// url, flagInput, convertLinks, pathRejects := mirror.GetMirrorUrl(inputs.args)
		mirror.DownloadPage(inputs.URL, inputs.RejectFlag, inputs.ConvertLinksFlag, inputs.ExcludeFlag)
		return
	}

	// If no file name is provided, derive it from the URL
	if inputs.File == "" && inputs.URL != "" {
		urlParts := strings.Split(inputs.URL, "/")
		inputs.File = urlParts[len(urlParts)-1]
	}

	// Handle the work-in-background flag
	if inputs.WorkInBackground {
		background.DownloadInBackground(inputs.File, inputs.URL, inputs.RateLimit)
		return
	}

	// Handle multiple file downloads from sourcefile
	if inputs.Sourcefile != "" {
		downloader.DownloadMultipleFiles(inputs.Sourcefile, inputs.File, inputs.RateLimit, inputs.Path)
		return
	}

	// Ensure URL is provided
	if inputs.URL == "" {
		fmt.Println("Error: URL not provided.")
		return
	}

	// Start downloading the file
	downloader.OneDownload(inputs.File, inputs.URL, inputs.RateLimit, inputs.Path)
}
