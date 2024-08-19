package downloader

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"downloaderex/internal/rateLimiter"
)

// Global map to keep track of processed URLs
var processedURLs = struct {
	sync.Mutex
	urls map[string]bool
}{
	urls: make(map[string]bool),
}

func DownloadMultipleFiles(filePath, outputFile, limit, directory string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var wg sync.WaitGroup
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue // Skip empty lines
		}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			AsyncDownload(outputFile, url, limit, directory)
		}(url)
	}
	wg.Wait()
}

func AsyncDownload(outputFileName, urlStr, limit, directory string) {
	// Check if the URL has already been processed
	processedURLs.Lock()
	if processed, exists := processedURLs.urls[urlStr]; exists && processed {
		processedURLs.Unlock()
		fmt.Printf("URL already processed: %s\n", urlStr)
		return
	}
	processedURLs.Unlock()

	// startTime := time.Now()
	// fmt.Printf("Start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	// Parse the URL to get the path components
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Create the necessary directories based on the URL path
	rootPath := ExpandPath(directory)
	pathComponents := strings.Split(strings.Trim(u.Path, "/"), "/")
	relativeDirPath := filepath.Join(pathComponents[:len(pathComponents)-1]...)
	fullDirPath := filepath.Join(rootPath, relativeDirPath)
	fileName := pathComponents[len(pathComponents)-1]

	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s\n", resp.Status)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	// contentLength := resp.ContentLength
	// fmt.Printf("Content size: %d bytes [~%.2fMB]\n", contentLength, float64(contentLength)/1024/1024)

	if outputFileName == "" {
		if fileName == "" || strings.HasSuffix(urlStr, "/") {
			fileName = "index.html"
		} else if contentType == "text/html" && !strings.HasSuffix(fileName, ".html") {
			fileName += ".html"
		}
		outputFileName = filepath.Join(fullDirPath, fileName)
	} else {
		if contentType == "text/html" && !strings.HasSuffix(outputFileName, ".html") {
			outputFileName += ".html"
		}
		outputFileName = filepath.Join(fullDirPath, outputFileName)
	}

	if fullDirPath != "" {
		if _, err := os.Stat(fullDirPath); os.IsNotExist(err) {
			err = os.MkdirAll(fullDirPath, 0o755)
			if err != nil {
				fmt.Println("Error creating path:", err)
				return
			}
		}
	}

	var out *os.File
	out, err = os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer out.Close()

	var reader io.Reader = resp.Body
	if limit != "" {
		reader = rateLimiter.NewRateLimitedReader(resp.Body, limit)
	}

	buffer := make([]byte, 32*1024) // 32 KB buffer size
	var downloaded int64
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
			downloaded += int64(n)
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Println() // Move to the next line after download completes

	// endTime := time.Now()
	fmt.Printf("Downloaded [%s]\n", urlStr)
	// fmt.Printf("Finished at %s\n", endTime.Format("2006-01-02 15:04:05"))

	// Mark the URL as processed
	processedURLs.Lock()
	processedURLs.urls[urlStr] = true
	processedURLs.Unlock()
}
