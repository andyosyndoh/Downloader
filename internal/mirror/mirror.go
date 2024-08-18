package mirror

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"downloaderex/internal/downloader"

	"golang.org/x/net/html"
)

var (
	visited = make(map[string]bool)
	mu      sync.Mutex
)

func extractDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

func DownloadPage(url, rejectTypes string) {
	domain, err := extractDomain(url)
	if err != nil {
		fmt.Println("Cold not extract domain name for:", url, "Error: ", err)
		return
	}
	// Check if URL has already been visited
	mu.Lock()
	if visited[url] {
		mu.Unlock()
		return
	}
	visited[url] = true
	mu.Unlock()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: status %s\n", resp.Status)
		return
	}
	// Read the response body and parse it
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Find and download assets
	var wg sync.WaitGroup
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			var link string
			for _, attr := range n.Attr {
				if (n.Data == "a" && attr.Key == "href") ||
					(n.Data == "img" && attr.Key == "src") ||
					(n.Data == "script" && attr.Key == "src") ||
					(n.Data == "link" && attr.Key == "href") {

					link = attr.Val
					if link != "" {
						wg.Add(1)
						go func(link string, tagName string) {
							defer wg.Done()
							baseURL := resolveURL(url, link)
							baseURLDomain, err := extractDomain(baseURL)
							if err != nil {
								fmt.Println("Cold not extract domain name for:", baseURLDomain, "Error: ", err)
								return
							}

							if baseURLDomain == domain {
								if tagName == "a" {
									DownloadPage(baseURL, rejectTypes) // Recursively download HTML pages
									downloadAsset(baseURL, domain, rejectTypes)
								} else {
									downloadAsset(baseURL, domain, rejectTypes) // Download assets like images, scripts, etc.
								}
							}
						}(link, n.Data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	// Wait for all downloads to complete
	wg.Wait()

	fmt.Println("Mirroring completed.")
}

func downloadAsset(fileURL, domain, rejectTypes string) {
	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	if isRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}

	fmt.Printf("Downloading: %s\n", fileURL)
	downloader.AsyncDownload("", fileURL, "", domain)
}

func resolveURL(base, rel string) string {
	if strings.HasPrefix(rel, "http") {
		return rel
	}

	if strings.HasPrefix(rel, "//") {
		protocol := "http:"
		if strings.HasPrefix(base, "https") {
			protocol = "https:"
		}
		return protocol + rel
	}

	if strings.HasPrefix(rel, "/") {
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel
	}
	if strings.HasPrefix(rel, "./") {
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel
	}
	if strings.HasPrefix(rel, "//") && strings.Contains(rel[2:], "/") {
		baseParts := strings.Split(base, "/")
		return baseParts[0] + "//" + baseParts[2] + rel[1:]
	}

	baseParts := strings.Split(base, "/")
	return baseParts[0] + "//" + baseParts[2] + "/" + rel
}

func isRejected(url, rejectTypes string) bool {
	if rejectTypes == "" {
		return false
	}

	rejectedTypes := strings.Split(rejectTypes, ",")
	for _, ext := range rejectedTypes {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}

func GetMirrorUrl(args []string) (string, string) {
	var url string
	var flagInput string
	for _, arg := range args {
		if strings.HasPrefix(arg, "http") {
			url = arg
		}
		if strings.HasPrefix(arg, "-R=") {
			flagInput = strings.TrimPrefix(arg, "-R=")
		}
	}
	return url, flagInput
}
