package mirror

import (
	"net/url"
	"os"
	"regexp"
	"strings"
)

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

func isRejectedPath(url, pathRejects string) bool {
	if pathRejects == "" {
		return false
	}

	rejects := strings.Split(pathRejects, ",")
	for _, path := range rejects {
		if path[0] != '/' {
			continue
		}
		if contains(url, path) {
			return true
		}
	}

	return false
}

func contains(str, substr string) bool {
	for i := 0; i < len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func extractAndHandleStyleURLs(styleContent, baseURL, domain, rejectTypes string) {
	re := regexp.MustCompile(`url\(['"]?([^'"()]+)['"]?\)`)
	matches := re.FindAllStringSubmatch(styleContent, -1)
	for _, match := range matches {
		if len(match) > 1 {
			assetURL := resolveURL(baseURL, match[1])
			downloadAsset(assetURL, domain, rejectTypes)
		}
	}
}

func extractDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}

// isValidAttribute checks if an HTML tag attribute is valid for processing
func isValidAttribute(tagName, attrKey string) bool {
	return (tagName == "a" && attrKey == "href") ||
		(tagName == "img" && attrKey == "src") ||
		(tagName == "script" && attrKey == "src") ||
		(tagName == "link" && attrKey == "href")
}
