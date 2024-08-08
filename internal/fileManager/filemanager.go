package fileManager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

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
