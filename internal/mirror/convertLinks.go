package mirror

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func convertLinks(htmlFilePath string) {
	htmlFilePath = removeHTTP(htmlFilePath)

	if !strings.HasSuffix(htmlFilePath, ".html") {
		return
	}
	// Open the HTML file for reading
	htmlFile, err := os.Open(htmlFilePath)
	if err != nil {
		fmt.Println("Error opening HTML file:", err)
		return
	}
	defer htmlFile.Close()

	// Read the HTML file content
	htmlData, err := ioutil.ReadAll(htmlFile)
	if err != nil {
		fmt.Println("Error reading HTML file:", err)
		return
	}

	// Parse the HTML content
	doc, err := html.Parse(strings.NewReader(string(htmlData)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Modify the document by converting external links to local paths
	modifyLinks(doc, path.Dir(htmlFilePath))

	// Convert the modified HTML back to string
	var modifiedHTML strings.Builder
	err = html.Render(&modifiedHTML, doc)
	if err != nil {
		fmt.Println("Error rendering modified HTML:", err)
		return
	}

	// Save the modified HTML back to the file
	err = ioutil.WriteFile(htmlFilePath, []byte(modifiedHTML.String()), 0o644)
	if err != nil {
		fmt.Println("Error writing modified HTML file:", err)
		return
	}

	fmt.Println("Links converted for offline viewing in", htmlFilePath)
}

func modifyLinks(n *html.Node, basePath string) {
	if n.Type == html.ElementNode {
		for i, attr := range n.Attr {
			if attr.Key == "href" || attr.Key == "src" {
				n.Attr[i].Val = getLocalPath(attr.Val)
			} else if attr.Key == "style" {
				n.Attr[i].Val = convertCSSURLs(attr.Val)
			}
		}

		if n.Data == "style" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
			n.FirstChild.Data = convertCSSURLs(n.FirstChild.Data)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		modifyLinks(c, basePath)
	}
}

func convertCSSURLs(cssContent string) string {
	re := regexp.MustCompile(`url\(([^)]+)\)`)
	return re.ReplaceAllStringFunc(cssContent, func(match string) string {
		url := strings.Trim(match[4:len(match)-1], "'\"")
		localPath := getLocalPath(url)
		return fmt.Sprintf("url('%s')", localPath)
	})
}

func getLocalPath(originalURL string) string {
	if strings.HasPrefix(originalURL, "http") || strings.HasPrefix(originalURL, "//") {
		parsedURL, err := url.Parse(originalURL)
		if err != nil {
			return originalURL
		}
		return path.Join(parsedURL.Host, parsedURL.Path)
	} else if strings.HasPrefix(originalURL, "/") {
		return path.Join(".", originalURL)
	}
	return originalURL
}

// removeHTTP removes the http:// or https:// prefix from the URL.
func removeHTTP(url string) string {
	// Regular expression to match the protocol (http or https) at the start of the URL
	re := regexp.MustCompile(`^https?://`)

	// Remove http or https from the URL
	modifiedURL := re.ReplaceAllString(url, "")

	// Check if the URL is a base URL (i.e., domain only without a path)
	// This regex checks if the modified URL is something like "example.com/"
	isBaseURL := regexp.MustCompile(`^[^/]+/?$`).MatchString(modifiedURL)

	// If the URL is a base URL, append "index.html" if it's not already present
	if isBaseURL {
		if strings.HasSuffix(modifiedURL, "/") {
			modifiedURL += "index.html"
		} else {
			modifiedURL += "/index.html"
		}
	}

	return modifiedURL
}

func IsFolder(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Check if the path is a directory
	return info.IsDir()
}
