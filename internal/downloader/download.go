package downloader

import (
	"downloaderex/internal/flags"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func OneDownload(file, url string) {
	fileURL := url
	startTime := time.Now()
	fmt.Printf("Start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s\n", resp.Status)
		return
	}
	fmt.Printf("Sending request, awaiting response... status %s\n", resp.Status)

	contentLength := resp.ContentLength
	fmt.Printf("Content size: %d bytes [~%.2fMB]\n", contentLength, float64(contentLength)/1024/1024)
	outputFile := ""
	if flags.OutputFileFlag(os.Args[1:]) {
		outputFile = flags.GetFlagInput(os.Args[1])
	} else {
		urlParts := strings.Split(fileURL, "/")
		fileName := urlParts[len(urlParts)-1]
		outputFile = "./" + fileName
	}
	fmt.Printf("Saving file to: %s\n", outputFile)

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64
	startDownload := time.Now()

	fmt.Print("Downloading... ")
	for {
		n, err := resp.Body.Read(buffer)
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

			// Calculate and display the progress
			progress := float64(downloaded) / float64(contentLength) * 100
			speed := float64(downloaded) / time.Since(startDownload).Seconds()
			timeRemaining := time.Duration(float64(contentLength-downloaded)/speed) * time.Second

			// Update the same line with progress
			fmt.Printf("\r%.2f KiB / %.2f KiB [", float64(downloaded)/1024, float64(contentLength)/1024)
			for i := 0; i < 100; i++ {
				if i < int(progress) {
					fmt.Print("=")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Printf("] %.2f%% %.2f KiB/s %s", progress, speed/1024, timeRemaining.String())

		}

		if downloaded >= contentLength {
			break
		}
	}

	fmt.Println() // Move to the next line after download completes

	endTime := time.Now()
	fmt.Printf("Downloaded [%s]\n", fileURL)
	fmt.Printf("Finished at %s\n", endTime.Format("2006-01-02 15:04:05"))
}
