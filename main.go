package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"downloaderex/internal/background"
	"downloaderex/internal/downloader"
	"downloaderex/internal/fileManager"
	"downloaderex/internal/mirror"
)

// Inputs struct to store command-line arguments and parsed values
type Inputs struct {
	args             []string // to store command line arguments
	url              string   // URL to download
	file             string   // file name for output
	rateLimit        string   // rate limit for downloading (e.g., 200k, 2M)
	path             string   // path to save the file
	sourcefile       string   // file containing URLs to download
	workInBackground bool     // whether to download in the background
	log              bool     // whether to log the output
	mirroring        bool     // whether to mirror a website
	rejectFlag       string
	excludeFlag      string
	convertLinksFlag bool
}

func main() {
	// Check if arguments are provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <URL> [options]")
		return
	}

	// Create an instance of Inputs and set command-line arguments
	inputs := Inputs{
		args: os.Args[1:],
	}

	inputs = ParseArgs()

	// Mirror website handling
	if inputs.mirroring {
		url, flagInput, convertLinks, pathRejects := mirror.GetMirrorUrl(inputs.args)
		mirror.DownloadPage(url, flagInput, convertLinks, pathRejects)
		return
	}

	// If no file name is provided, derive it from the URL
	if inputs.file == "" && inputs.url != "" {
		urlParts := strings.Split(inputs.url, "/")
		inputs.file = urlParts[len(urlParts)-1]
	}

	// Handle the work-in-background flag
	if inputs.workInBackground {
		background.DownloadInBackground(inputs.file, inputs.url, inputs.rateLimit)
		return
	}

	// Handle the log flag
	if inputs.log {
		fileManager.Logger(inputs.file, inputs.url, inputs.rateLimit)
		return
	}

	// Handle multiple file downloads from sourcefile
	if inputs.sourcefile != "" {
		downloader.DownloadMultipleFiles(inputs.sourcefile, inputs.file, inputs.rateLimit, inputs.path)
		return
	}

	// Ensure URL is provided
	if inputs.url == "" {
		fmt.Println("Error: URL not provided.")
		return
	}

	// Start downloading the file
	downloader.OneDownload(inputs.file, inputs.url, inputs.rateLimit, inputs.path)
}

func ParseArgs() Inputs {
	input := &Inputs{}
	// Parse the command-line arguments
	input.workInBackground = *flag.Bool("B", false, "Download file in the background")
	input.file = *flag.String("O", "", "Save file under a different name")
	input.path = *flag.String("P", "", "Path to save the file")
	input.rateLimit = *flag.String("rate-limit", "", "Limit download speed (e.g., 200k, 2M)")
	input.sourcefile = *flag.String("i", "", "File containing links to download asynchronously")
	input.mirroring = *flag.Bool("mirror", false, "Mirror a website for offline use")
	input.rejectFlag = *flag.String("R", "", "File suffixes to avoid downloading during retrieval")
	input.excludeFlag = *flag.String("X", "", "Directories to exclude from download")
	input.convertLinksFlag = *flag.Bool("convert-links", false, "Convert links for offline viewing")

	// Parse flags
	flag.Parse()

	// After parsing flags, the remaining arguments should contain the URL(s)
	remainingArgs := flag.Args()
	if len(remainingArgs) == 0 {
		fmt.Println("no URL provided")
		os.Exit(0)
	}
	input.url = remainingArgs[0]

	// Validate the URL
	err := validateURL(input.url)
	if err != nil {
		fmt.Println("no URL provided")
		os.Exit(0)
	}

	return *input

	// for _, arg := range inputs.args {
	// 	if strings.HasPrefix(arg, "--rate-limit=") {
	// 		if inputs.rateLimit != "" {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.rateLimit = arg[len("--rate-limit="):]
	// 	} else if strings.HasPrefix(arg, "-O=") {
	// 		if inputs.file != "" {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.file = flags.GetFlagInput(arg)
	// 	} else if strings.HasPrefix(arg, "-P=") {
	// 		if inputs.path != "" {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.path = flags.GetFlagInput(arg)
	// 	} else if strings.HasPrefix(arg, "http") {
	// 		if inputs.url != "" {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.url = arg
	// 	} else if strings.HasPrefix(arg, "-B") {
	// 		if inputs.workInBackground {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.workInBackground = true
	// 	} else if arg == "-b" {
	// 		if inputs.log {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.log = true
	// 	} else if strings.HasPrefix(arg, "-i=") {
	// 		if inputs.sourcefile != "" {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.sourcefile = flags.GetFlagInput(arg)
	// 	} else if strings.HasPrefix(arg, "--mirror") {
	// 		if inputs.mirroring {
	// 			fmt.Printf("Error: Repeated argument '%v'\n", arg)
	// 			os.Exit(0)
	// 		}
	// 		inputs.mirroring = true
	// 		break
	// 	}
	// }
}

func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}
