package flags

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Inputs struct with exported fields (Uppercase names)
type Inputs struct {
	URL              string
	File             string
	RateLimit        string
	Path             string
	Sourcefile       string
	WorkInBackground bool
	Mirroring        bool // Capitalized "Mirroring"
	RejectFlag       string
	ExcludeFlag      string
	ConvertLinksFlag bool
}

func ParseArgs() Inputs {
	input := &Inputs{}
	mirrorMode := false // Flag to track if --mirror is set
	track := false

	// Iterate over the command-line arguments manually
	for _, arg := range os.Args[1:] {
		// Enforce flags with the '=' sign
		if strings.HasPrefix(arg, "-O=") {
			input.File = arg[len("-O="):] // Capture the file name
		} else if strings.HasPrefix(arg, "-P=") {
			input.Path = arg[len("-P="):] // Capture the path
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			input.RateLimit = arg[len("--rate-limit="):] // Capture the rate limit
		} else if strings.HasPrefix(arg, "--mirror") {
			input.Mirroring = true // Enable mirroring
			mirrorMode = true      // Track mirror mode
		} else if strings.HasPrefix(arg, "--convert-links") {
			if !mirrorMode {
				fmt.Println("Error: --convert-links can only be used with --mirror.")
				os.Exit(1)
			}
			input.ConvertLinksFlag = true // Enable link conversion
		} else if strings.HasPrefix(arg, "-R=") || strings.HasPrefix(arg, "--reject=") {
			if !mirrorMode {
				fmt.Println("Error: --reject can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-R=") {
				input.RejectFlag = arg[len("-R="):] // Capture reject flag for -R
			} else {
				input.RejectFlag = arg[len("--reject="):] // Capture reject flag for --reject
			}
		} else if strings.HasPrefix(arg, "-X=") || strings.HasPrefix(arg, "--exclude=") {
			if !mirrorMode {
				fmt.Println("Error: --exclude can only be used with --mirror.")
				os.Exit(1)
			}
			if strings.HasPrefix(arg, "-X=") {
				input.ExcludeFlag = arg[len("-X="):] // Capture exclude flag for -X
			} else {
				input.ExcludeFlag = arg[len("--exclude="):] // Capture exclude flag for --exclude
			}
		} else if strings.HasPrefix(arg, "-B") {
			input.WorkInBackground = true // Enable background downloading
		} else if strings.HasPrefix(arg, "-i=") {
			input.Sourcefile = arg[len("-i="):] // Capture source file
			track = true
		} else if strings.HasPrefix(arg, "http") {
			// This must be the URL
			input.URL = arg
		} else {
			fmt.Printf("Error: Unrecognized argument '%s'\n", arg)
			os.Exit(1)
		}
	}
	if input.RateLimit != "" {
		if strings.ToLower(string(input.RateLimit[len(input.RateLimit)-1])) != "k" &&
			strings.ToLower(string(input.RateLimit[len(input.RateLimit)-1])) != "m" {
			fmt.Println("Invalid RateLimit")
			os.Exit(1)
		}
	}

	// Check for invalid flag combinations if --mirror is provided
	if input.Mirroring {
		// Only allow --convert-links, --reject, and --exclude with --mirror
		if input.File != "" || input.Path != "" || input.RateLimit != "" || input.Sourcefile != "" || input.WorkInBackground {
			fmt.Println("Error: --mirror can only be used with --convert-links, --reject, --exclude, and a URL. No other flags are allowed.")
			os.Exit(1)
		}
	} else {
		// If --mirror is not provided, reject the use of --convert-links, --reject, and --exclude
		if input.ConvertLinksFlag || input.RejectFlag != "" || input.ExcludeFlag != "" {
			fmt.Println("Error: --convert-links, --reject, and --exclude can only be used with --mirror.")
			os.Exit(1)
		}
	}

	// Ensure URL is provided
	if input.URL == "" && !track {
		fmt.Println("Error: URL not provided.")
		os.Exit(1)

		// Validate the URL
		err := validateURL(input.URL)
		if err != nil {
			fmt.Println("Error: invalid URL provided")
			os.Exit(1)
		}
	}

	return *input
}

func validateURL(link string) error {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	return nil
}
