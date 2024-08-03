package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println("Usage: go run . <URL>")
	// 	return
	// }
	if len(os.Args) == 2 {
		file := ""
		url := os.Args[1]
		OneDownload(file, url)
	}
	if len(os.Args) == 3 && os.Args[1] == "-B" {
		Logger()
	} else if len(os.Args) == 3 && os.Args[1][:3] == "-O=" {
		file := os.Args[1][3:]
		url := os.Args[2]
		OneDownload(file, url)
	}

	fmt.Println("hey")
}

func logToFile(logFile *os.File, message string) {
	logFile.WriteString(message + "\n")
}

func Logger() {
	// Open log file
	logFile, err := os.OpenFile("wget-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	// Get the URL from the command-line argument
	fileURL := os.Args[2]

	// Record the start time
	startTime := time.Now()
	logToFile(logFile, fmt.Sprintf("Start at %s", startTime.Format("2006-01-02 15:04:05")))

	// Send a GET request to the specified URL
	resp, err := http.Get(fileURL)
	if err != nil {
		logToFile(logFile, fmt.Sprintf("Error downloading file: %v", err))
		return
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		logToFile(logFile, fmt.Sprintf("Error: status %s", resp.Status))
		return
	}
	logToFile(logFile, fmt.Sprintf("Sending request, awaiting response... status %s", resp.Status))

	// Get the content length
	contentLength := resp.ContentLength
	logToFile(logFile, fmt.Sprintf("Content size: %d [~%.2fMB]", contentLength, float64(contentLength)/1024/1024))

	outputFile := ""

	urlParts := strings.Split(fileURL, "/")
	fileName := urlParts[len(urlParts)-1]
	outputFile = "./" + fileName

	logToFile(logFile, fmt.Sprintf("Saving file to: %s", outputFile))

	// Create a new file to save the downloaded content
	out, err := os.Create(outputFile)
	if err != nil {
		logToFile(logFile, fmt.Sprintf("Error creating file: %v", err))
		return
	}
	defer out.Close()

	// Buffer to hold chunks of data being read
	buffer := make([]byte, 32*1024) // 32 KB buffer size

	for {
		// Read a chunk of data from the response body
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			logToFile(logFile, fmt.Sprintf("Error reading response body: %v", err))
			return
		}
		if n == 0 {
			break
		}

		// Write the chunk to the file
		if _, err := out.Write(buffer[:n]); err != nil {
			logToFile(logFile, fmt.Sprintf("Error writing to file: %v", err))
			return
		}

	}

	// Record the end time
	endTime := time.Now()
	logToFile(logFile, fmt.Sprintf("Downloaded [%s]", fileURL))
	logToFile(logFile, fmt.Sprintf("Finished at %s", endTime.Format("2006-01-02 15:04:05")))
}

func OneDownload(file, url string) {
	// Get the URL from the command-line argument
	fileURL := url

	// Record the start time
	startTime := time.Now()
	fmt.Printf("Start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	// Send a GET request to the specified URL
	resp, err := http.Get(fileURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s\n", resp.Status)
		return
	}
	fmt.Printf("Sending request, awaiting response... status %s\n", resp.Status)

	// Get the content length
	contentLength := resp.ContentLength
	fmt.Printf("Content size: %d [~%.2fMB]\n", contentLength, float64(contentLength)/1024/1024)

	outputFile := ""
	if os.Args[1][:3] == "-O=" {
		outputFile = file
	} else {
		urlParts := strings.Split(fileURL, "/")
		fileName := urlParts[len(urlParts)-1]
		outputFile = "./" + fileName
	}
	fmt.Printf("Saving file to: %s\n", outputFile)

	// Create a new file to save the downloaded content
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	// Buffer to hold chunks of data being read
	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64

	// Start downloading with progress reporting
	startDownload := time.Now()
	for {
		// Read a chunk of data from the response body
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading response body:", err)
			return
		}
		if n == 0 {
			break
		}

		// Write the chunk to the file
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

		fmt.Printf("\r%.2f MiB / %.2f MiB [%.2f%%] %.2f MiB/s %s",
			float64(downloaded)/1024/1024,
			float64(contentLength)/1024/1024,
			progress,
			speed/1024/1024,
			timeRemaining.String())
	}

	fmt.Println()

	// Record the end time
	endTime := time.Now()
	fmt.Printf("Downloaded [%s]\n", fileURL)
	fmt.Printf("Finished at %s\n", endTime.Format("2006-01-02 15:04:05"))
}
