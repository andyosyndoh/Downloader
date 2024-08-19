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

// Global variables to keep track of visited URLs and synchronization
var (
	visited = make(map[string]bool)
	mu      sync.Mutex // Mutex for thread-safe operations on 'visited'
)

// DownloadPage downloads a page and its assets, recursively visiting links
func DownloadPage(url, rejectTypes string) {
	domain, err := extractDomain(url)
	if err != nil {
		fmt.Println("Cold not extract domain name for:", url, "Error: ", err)
		return
	}
	// Check if URL has already been visited
	if !shouldDownload(url) {
		return
	}

	// Fetch and get the HTML of the page
	doc, err := fetchAndParsePage(url)
	if err != nil {
		fmt.Println("Error fetching or parsing page:", err)
		return
	}

	// Function to handle links and assets found on the page
	handleLink := func(link, tagName string) {
		baseURL := resolveURL(url, link)
		baseURLDomain, err := extractDomain(baseURL)
		if err != nil {
			fmt.Println("Could not extract domain name for:", baseURLDomain, "Error:", err)
			return
		}
		if baseURLDomain == domain {
			if tagName == "a" {
				DownloadPage(baseURL, rejectTypes)
			}
			downloadAsset(baseURL, domain, rejectTypes)
		}
	}

	var wg sync.WaitGroup
	var processNode func(n *html.Node)
	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if isValidAttribute(n.Data, attr.Key) {
					link := attr.Val
					if link != "" {
						wg.Add(1)
						go func(link, tagName string) {
							defer wg.Done()
							handleLink(link, tagName)
						}(link, n.Data)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c) // Recursively process child nodes
		}
	}
	processNode(doc)

	// Wait for all downloads to complete
	wg.Wait()

	fmt.Println("Mirroring completed.")
}

func extractDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// shouldDownload determines if a URL should be downloaded based on whether it has been visited
func shouldDownload(url string) bool {
	mu.Lock()
	defer mu.Unlock()
	if visited[url] {
		return false
	}
	visited[url] = true
	return true
}

// fetchAndParsePage fetches the content of the URL and parses it as HTML
func fetchAndParsePage(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status %s", resp.Status)
	}

	return html.Parse(resp.Body)
}

// isValidAttribute checks if an HTML tag attribute is valid for processing
func isValidAttribute(tagName, attrKey string) bool {
	return (tagName == "a" && attrKey == "href") ||
		(tagName == "img" && attrKey == "src") ||
		(tagName == "script" && attrKey == "src") ||
		(tagName == "link" && attrKey == "href")
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
