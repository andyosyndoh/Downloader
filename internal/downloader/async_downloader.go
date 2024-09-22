package downloader

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"downloaderex/internal/rateLimiter"
)

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
func AsyncDownload(outputFileName, url, limit, directory string) {
	path := ExpandPath(directory)
	// startTime := time.Now()
	// fmt.Printf("Start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	req.Header.Set("Referer", "http://ipv4.download.thinkbroadband.com/") // Change to a referer appropriate for the URL

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s url: [%s]\n", resp.Status, url)
		return
	}

	if outputFileName == "" {
		urlParts := strings.Split(url, "/")
		fileName := urlParts[len(urlParts)-1]
		outputFileName = filepath.Join(path, fileName)
	} else {
		outputFileName = filepath.Join(path, outputFileName)
	}

	if path != "" {
		err = os.MkdirAll(path, 0o755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
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

	// endTime := time.Now()
	fmt.Printf("\033[32mDownloaded\033[0m [%s]\n", url)
	// fmt.Printf("Finished at %s\n", endTime.Format("2006-01-02 15:04:05"))
}
