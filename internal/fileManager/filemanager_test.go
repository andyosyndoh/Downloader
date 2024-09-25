package fileManager

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Test cases
	tests := []struct {
		name      string
		file      string
		url       string
		rateLimit string
		wantFile  string
		wantError bool
	}{
		{
			name:      "Valid URL, default file name",
			file:      "", 
			url:       "https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg",
			rateLimit: "",
			wantFile:  "EMtmPFLWkAA8CIS.jpg",
			wantError: false,
		},
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Chdir("../.."); err != nil {
				t.Fatalf("Failed to change working directory: %v", err)
			}
			// Ensure the working directory is reverted back after the test
			t.Cleanup(func() {
				if err := os.Chdir(originalWD); err != nil {
					t.Fatalf("Failed to revert working directory: %v", err)
				}
			})
			Logger(tt.file, tt.url, tt.rateLimit)

			// Wait a moment for the download to complete
			time.Sleep(5 * time.Second)

			outputPath := filepath.Join("./", tt.wantFile)

			// Check if the file exists
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not downloaded", outputPath)
			} else {
				if err := os.Remove(outputPath); err != nil {
					t.Errorf("Failed to remove downloaded file: %v", err)
				}
			}
		})
	}
}
