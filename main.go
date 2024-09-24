package main

import (
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
	// args             []string // to store command line arguments
	url              string // URL to download
	file             string // file name for output
	rateLimit        string // rate limit for downloading (e.g., 200k, 2M)
	path             string // path to save the file
	sourcefile       string // file containing URLs to download
	workInBackground bool   // whether to download in the background
	log              bool   // whether to log the output
	mirroring        bool   // whether to mirror a website
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
	// inputs := Inputs{
	// 	args: os.Args[1:],
	// }

	inputs := ParseArgs()

	// Mirror website handling
	if inputs.mirroring {
		// url, flagInput, convertLinks, pathRejects := mirror.GetMirrorUrl(inputs.args)
		mirror.DownloadPage(inputs.url, inputs.rejectFlag, inputs.convertLinksFlag, inputs.excludeFlag)
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
	mirrorMode := false // Flag to track if --mirror is set

	// Iterate over the command-line arguments manually
	for _, arg := range os.Args[1:] {
		// Enforce flags with the '=' sign
		if strings.HasPrefix(arg, "-O=") {
			input.file = arg[len("-O="):] // Capture the file name
		} else if strings.HasPrefix(arg, "-P=") {
			input.path = arg[len("-P="):] // Capture the path
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			input.rateLimit = arg[len("--rate-limit="):] // Capture the rate limit
		} else if strings.HasPrefix(arg, "--mirror") {
			input.mirroring = true // Enable mirroring
			mirrorMode = true      // Track mirror mode
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				fmt.Println("Error: --convert-links can only be used with --mirror.")
				os.Exit(1)
			}
			input.convertLinksFlag = true // Enable link conversion
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				fmt.Println("Error: --reject can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-R=") {
				input.rejectFlag = arg[len("-R="):] // Capture reject flag for -R
			} else {
				input.rejectFlag = arg[len("--reject="):] // Capture reject flag for --reject
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				fmt.Println("Error: --exclude can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-X=") {
				input.excludeFlag = arg[len("-X="):] // Capture exclude flag for -X
			} else {
				input.excludeFlag = arg[len("--exclude="):] // Capture exclude flag for --exclude
			}
		} else if strings.HasPrefix(arg, "-B") {
			input.workInBackground = true // Enable background downloading
		} else if strings.HasPrefix(arg, "-i=") {
			input.sourcefile = arg[len("-i="):] // Capture source file
		} else if strings.HasPrefix(arg, "http") {
			// This must be the URL
			input.url = arg
		} else {
			fmt.Printf("Error: Unrecognized argument '%s'\n", arg)
			os.Exit(1)
		}
	}

	// Check for invalid flag combinations if --mirror is provided
	if input.mirroring {
		// Only allow --convert-links, --reject, and --exclude with --mirror
		if input.file != "" || input.path != "" || input.rateLimit != "" || input.sourcefile != "" || input.workInBackground {
			fmt.Println("Error: --mirror can only be used with --convert-links, --reject, --exclude, and a URL. No other flags are allowed.")
			os.Exit(1)
		}
	} else {
		// If --mirror is not provided, reject the use of --convert-links, --reject, and --exclude
		if input.convertLinksFlag || input.rejectFlag != "" || input.excludeFlag != "" {
			fmt.Println("Error: --convert-links, --reject, and --exclude can only be used with --mirror.")
			os.Exit(1)
		}
	}

	// Ensure URL is provided
	if input.url == "" {
		fmt.Println("Error: URL not provided.")
		os.Exit(1)
	}
	// Validate the URL
	err := validateURL(input.url)
	if err != nil {
		fmt.Println("Error: invalid URL provided")
		os.Exit(1)
	}

	// Debugging prints to confirm that the flags are captured correctly
	// fmt.Println("Flags captured:")
	// fmt.Printf("Work in background: %v\n", input.workInBackground)
	// fmt.Printf("File name: %v\n", input.file)
	// fmt.Printf("Path: %v\n", input.path)
	// fmt.Printf("Rate limit: %v\n", input.rateLimit)
	// fmt.Printf("Source file: %v\n", input.sourcefile)
	// fmt.Printf("Mirroring: %v\n", input.mirroring)
	// fmt.Printf("Reject flag: %v\n", input.rejectFlag)
	// fmt.Printf("Exclude flag: %v\n", input.excludeFlag)
	// fmt.Printf("Convert links: %v\n", input.convertLinksFlag)
	// fmt.Printf("URL: %v\n", input.url)

	return *input
}

func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}
