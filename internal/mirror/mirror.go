package mirror

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"downloaderex/internal/downloader"

	"golang.org/x/net/html"
)

func extractDomain(url string) string {
	// Logic to extract domain from the URL
	// This can be a simple implementation
	parts := strings.Split(strings.TrimPrefix(url, "https://"), "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func DownloadPage(url, rejectTypes string) {
	domain := extractDomain(url)

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
			if n.Data == "a" || n.Data == "link" || n.Data == "img" {
				for _, attr := range n.Attr {
					// Use the base URL derived from the current page
					if (n.Data == "a" || n.Data == "link") && attr.Key == "href" {
						wg.Add(1) // Increment the WaitGroup counter
						go func(link string) {
							defer wg.Done() // Decrement the counter when done
							// Resolve the base URL for links
							baseURL := resolveURL(url, link)
							downloadAsset(baseURL, domain, rejectTypes)
						}(attr.Val) // Pass the URL
					}
					if n.Data == "img" && attr.Key == "src" {
						wg.Add(1) // Increment the WaitGroup counter
						go func(src string) {
							defer wg.Done() // Decrement the counter when done
							// Resolve the base URL for images
							baseURL := resolveURL(url, src)
							downloadAsset(baseURL, domain, rejectTypes)
						}(attr.Val) // Pass the URL
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
	// Check if the URL is valid after resolving
	if fileURL == "" || !strings.HasPrefix(fileURL, "http") {
		fmt.Printf("Invalid URL: %s\n", fileURL)
		return
	}

	// Check if file type is rejected
	if isRejected(fileURL, rejectTypes) {
		fmt.Printf("Skipping rejected file: %s\n", fileURL)
		return
	}

	fmt.Printf("Downloading: %s\n", fileURL)

	// Call your OneDownload function here
	downloader.AsyncDownload("", fileURL, "", domain)
}

func resolveURL(base, rel string) string {
	// If the relative URL starts with a scheme, return it as is.
	if strings.HasPrefix(rel, "http") {
		return rel
	}

	// Handle protocol-relative URLs (starting with //)
	if strings.HasPrefix(rel, "//") && !strings.Contains(rel[2:], "/") {
		// Determine the protocol of the base URL

		protocol := "http:"
		if strings.HasPrefix(base, "https") {
			protocol = "https:"
		}
		return protocol + rel
	}

	// Handle relative paths
	if strings.HasPrefix(rel, "/") {
		// For absolute paths, return base URL without the last segment
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel
	}
	if strings.HasPrefix(rel, "./") {
		// For absolute paths, return base URL without the last segment
		return strings.Join(strings.Split(base, "/")[:3], "/") + rel
	}
	// Handle paths that start with // but are not protocol-relative (e.g., //css/style.css)
	if strings.HasPrefix(rel, "//") && strings.Contains(rel[2:], "/") {

		baseParts := strings.Split(base, "/")
		return baseParts[0] + "//" + baseParts[2] + rel[1:] // Treat it as a root-relative path
	}

	// Remove the last part of the base URL
	baseParts := strings.Split(base, "/")
	// baseParts = baseParts[:len(baseParts)-1] // remove last part

	// Append the relative URL
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
