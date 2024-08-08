package main

import (
	"downloaderex/internal/downloader"
	"downloaderex/internal/fileManager"
	"fmt"
	"os"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println("Usage: go run . <URL>")
	// 	return
	// }
	if len(os.Args) == 2 {
		file := ""
		url := os.Args[1]
		downloader.OneDownload(file, url)
	}
	if len(os.Args) == 3 && os.Args[1] == "-B" {
		fileManager.Logger()
		fmt.Println("Output will be written to \"wget-log\"")
	} else if len(os.Args) == 3 && os.Args[1][:3] == "-O=" {
		file := os.Args[1][3:]
		url := os.Args[2]
		downloader.OneDownload(file, url)

	}

}
