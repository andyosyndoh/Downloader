package mirror

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Global variables to keep track of visited URLs and synchronization
var (
	visitedPages  = make(map[string]bool)
	visitedAssets = make(map[string]bool)
	muPages       sync.Mutex
	muAssets      sync.Mutex
	semaphore     = make(chan struct{}, 50)
	count         int
)

// DownloadPage downloads a page and its assets, recursively visiting links
func DownloadPage(url, rejectTypes string) {
	domain, err := extractDomain(url)
	if err != nil {
		fmt.Println("Could not extract domain name for:", url, "Error:", err)
		return
	}

	muPages.Lock()
	if visitedPages[url] {
		muPages.Unlock()
		return
	}
	visitedPages[url] = true
	muPages.Unlock()

	// Fetch and get the HTML of the page
	doc, err := fetchAndParsePage(url)
	if err != nil {
		fmt.Println("Error fetching or parsing page:", err)
		return
	}

	// Function to handle links and assets found on the page
	handleLink := func(link, tagName string) {
		semaphore <- struct{}{}        // Acquire a spot in the semaphore
		defer func() { <-semaphore }() // Release the spot

		baseURL := resolveURL(url, link)

		baseURLDomain, err := extractDomain(baseURL)
		if err != nil {
			fmt.Println("Could not extract domain name for:", baseURLDomain, "Error:", err)
			return
		}

		if baseURLDomain == domain {
			if tagName == "a" {
				// Check if the baseURL is the root or equivalent to index.html
				if strings.HasSuffix(baseURL, "/") || strings.HasSuffix(baseURL, "/index.html") {
					// Ensure index.html is downloaded first
					if !visitedPages["http://"+baseURLDomain+"/index.html"] {
						DownloadPage("http://"+baseURLDomain+"/index.html", rejectTypes)
					}
					// Recursively process other pages
					DownloadPage(baseURL, rejectTypes)
				} else {
					// Process other pages as usual
					DownloadPage(baseURL, rejectTypes)
				}
			}
			// Download assets, regardless of index.html processing
			downloadAsset(baseURL, domain, rejectTypes)
		}
	}
	var wg sync.WaitGroup
	var processNode func(n *html.Node)
	processedPages := make(map[string]bool)

	processNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if isValidAttribute(n.Data, attr.Key) {
					link := attr.Val
					if link != "" {
						baseURL := resolveURL(url, link)
						baseURLDomain, err := extractDomain(baseURL)
						if err != nil {
							fmt.Println("Could not extract domain name for:", baseURLDomain, "Error:", err)
							continue
						}

						// Process index.html first
						if (strings.HasSuffix(url, ".com/") || strings.HasSuffix(url, ".com/index.html")) && count == 0 {
							count++
							if !processedPages["http://"+baseURLDomain+"/index.html"] {
								wg.Add(1)
								go func(link, tagName string) {
									defer wg.Done()
									handleLink(link, tagName)
									processedPages[url] = true
								}("http://"+baseURLDomain+"/index.html", n.Data)
							}
						}

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

	// Start processing the document
	processNode(doc)

	// Wait for all goroutines to complete
	wg.Wait()

}

func extractDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
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

	// Remove fragment identifiers (anything starting with #)
	if fragmentIndex := strings.Index(rel, "#"); fragmentIndex != -1 {
		rel = rel[:fragmentIndex]
	}

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
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel[1:]
	}
	if strings.HasPrefix(rel, "//") && strings.Contains(rel[2:], "/") {
		baseParts := strings.Split(base, "/")
		return baseParts[0] + "//" + baseParts[2] + rel[1:]
	}

	baseParts := strings.Split(base, "/")
	return baseParts[0] + "//" + baseParts[2] + "/" + rel
}

func downloadAsset(fileURL, domain, rejectTypes string) {
	muAssets.Lock()
	if visitedAssets[fileURL] {
		muAssets.Unlock()
		return
	}
	visitedAssets[fileURL] = true
	muAssets.Unlock()

	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	if isRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}
	fmt.Printf("Downloading: %s\n", fileURL)
	MirrorAsyncDownload("", fileURL, "", domain)
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
