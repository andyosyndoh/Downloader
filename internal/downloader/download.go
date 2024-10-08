package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wiget/internal/background"
	"wiget/internal/rateLimiter"
)

func OneDownload(file, url, limit, directory string) {
	path := ExpandPath(directory)
	fileURL := url
	startTime := time.Now()
	toDisplay, err := background.LoadShowProgressState()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	resp, err := HttpRequest(fileURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s url: [%s]\n", resp.Status, url)
		return
	}
	fmt.Printf("sending request, awaiting response... status %s\n", resp.Status)

	contentLength := resp.ContentLength
	fmt.Printf("content size: %d bytes [~%.2fMB]\n", contentLength, float64(contentLength)/1000000)

	// Set the output file name
	var outputFile string
	if file == "" {
		urlParts := strings.Split(fileURL, "/")
		file = urlParts[len(urlParts)-1]
		outputFile = filepath.Join(path, file)
	} else {
		outputFile = filepath.Join(path, file)
	}
	// Create the path if it doesn't exist
	if path != "" {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			fmt.Println("Error creating path:", err)
			return
		}
	}
	temp := ""
	if file != "" && directory != "" {
		fmt.Printf("saving file to: %s%s\n", directory, file)
	} else if path == "" && file != "" {
		temp = "./"
		fmt.Printf("saving file to: %s%s\n", temp, file)
	} else {
		temp = "./"
		fmt.Printf("saving file to: %s%s\n", temp, file)
	}

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	var reader io.Reader
	if limit != "" {
		reader = rateLimiter.NewRateLimitedReader(resp.Body, limit) // Assuming rateLimiter is defined elsewhere
	} else {
		reader = resp.Body
	}

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64
	startDownload := time.Now()

	if toDisplay {
		fmt.Print("Downloading... ")
	}
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading response body:", err)
			return
		}

		if n > 0 {
			if _, err := out.Write(buffer[:n]); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			// Update the downloaded size
			downloaded += int64(n)

			if toDisplay {
				// Calculate and display the progress
				progress := float64(downloaded) / float64(contentLength) * 50
				speed := float64(downloaded) / time.Since(startDownload).Seconds()
				timeRemaining := time.Duration(float64(contentLength-downloaded)/speed) * time.Second

				// Update the same line with progress
				fmt.Printf("\r %.2f KiB / %.2f KiB [", float64(downloaded)/1024, float64(contentLength)/1024)
				for i := 0; i < 50; i++ {
					if i < int(progress) {
						fmt.Print("=")
					} else {
						fmt.Print(" ")
					}
				}
				fmt.Printf("] %.2f%% %.2f KiB/s %s", (float64(downloaded)*100)/float64(contentLength), speed/1024, timeRemaining.String())

			}

		}

		if downloaded >= contentLength {
			break
		}
	}
	if toDisplay {
		fmt.Println() // Move to the next line after download completes
		fmt.Println()
	}

	endTime := time.Now()
	fmt.Printf("Downloaded [%s]\n", fileURL)
	fmt.Printf("finished at %s\n", endTime.Format("2006-01-02 15:04:05"))
	if !toDisplay {
		fmt.Println()
	}
}

// ExpandPath expands shorthand notations to full paths
func ExpandPath(path string) string {
	// 1. Expand `~` to the home directory
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error finding home directory:", err)
			return ""
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// 2. Expand environment variables like $HOME, $USER, etc.
	path = os.ExpandEnv(path)

	// 3. Convert relative paths (./ or ../) to absolute paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return ""
	}

	return absPath
}
