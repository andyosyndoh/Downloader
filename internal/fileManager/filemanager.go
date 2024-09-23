package fileManager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"downloaderex/internal/downloader"
	"downloaderex/internal/flags"
	"downloaderex/internal/rateLimiter"
)

func logToFile(logFile *os.File, message string) {
	logFile.WriteString(message + "\n")
}

func Logger(file, url, limit string) {
	// Open log file
	logFile, err := os.OpenFile("wget-log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()
	fileURL := url
	startTime := time.Now()

	logToFile(logFile, fmt.Sprintf("start at %s", startTime.Format("2006-01-02 15:04:05")))

	resp, err := downloader.HttpRequest(fileURL)
	if err != nil {
		logToFile(logFile, fmt.Sprintf("Error downloading file: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logToFile(logFile, fmt.Sprintf("Error: status %s\n", resp.Status))
		return
	}
	logToFile(logFile, fmt.Sprintf("sending request, awaiting response... status %s", resp.Status))

	contentLength := resp.ContentLength
	logToFile(logFile, fmt.Sprintf("content size: %d bytes [~%.2fMB]", contentLength, float64(contentLength)/1024/1024))

	outputFile := ""
	if flags.OutputFileFlag(os.Args[1:]) {
		outputFile = flags.GetFlagInput(os.Args[1])
	} else {
		urlParts := strings.Split(fileURL, "/")
		fileName := urlParts[len(urlParts)-1]
		outputFile = "./" + fileName
	}
	logToFile(logFile, fmt.Sprintf("saving file to: %s", outputFile))

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Sprintln("Error creating file:", err)
		return
	}
	defer out.Close()

	var reader io.Reader

	if limit != "" {
		reader = rateLimiter.NewRateLimitedReader(resp.Body, limit)
	} else {
		reader = resp.Body
	}

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			logToFile(logFile, fmt.Sprintf("Error reading response body: %v", err))
			return
		}

		if n > 0 {
			if _, err := out.Write(buffer[:n]); err != nil {
				logToFile(logFile, fmt.Sprintf("Error writing to file: %v", err))
				return
			}
			// Update the downloaded size
			downloaded += int64(n)
		}

		if downloaded >= contentLength {
			break
		}
	}

	fmt.Sprintln() // Move to the next line after download completes

	endTime := time.Now()
	logToFile(logFile, fmt.Sprintf("Downloaded [%s]", fileURL))
	logToFile(logFile, fmt.Sprintf("finished at %s\n", endTime.Format("2006-01-02 15:04:05")))
}
