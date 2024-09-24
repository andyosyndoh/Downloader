package downloader

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create a mock test file with given URLs.
func createMockFile(filePath string, urls []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, url := range urls {
		_, err := writer.WriteString(url + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// Helper function to clean up the mock test files.
func cleanupMockFiles(files []string) {
	for _, file := range files {
		os.Remove(file)
	}
}

// Silence the terminal output during tests.
func silenceOutput() (func(), func()) {
	// Save the original stdout and stderr
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Create a nil writer to silence output
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// Restore original output streams
	restore := func() {
		w.Close()
		os.Stdout = origStdout
		os.Stderr = origStderr
	}
	return restore, func() {
		r.Close()
	}
}

func TestDownloadMultipleFiles(t *testing.T) {
	type args struct {
		filePath   string
		outputFile string
		limit      string
		directory  string
	}
	tests := []struct {
		name       string
		args       args
		mockURLs   []string
		expectFail bool
	}{
		{
			name: "Valid URLs with no limit",
			args: args{
				filePath:   "./test_urls.txt",
				outputFile: "output.txt",
				limit:      "",
				directory:  "./",
			},
			mockURLs:   []string{"https://example.com/file1.txt", "https://example.com/file2.txt"},
			expectFail: true,
		},
		{
			name: "Valid URLs with limit",
			args: args{
				filePath:   "./test_urls_with_limit.txt",
				outputFile: "output_with_limit.txt",
				limit:      "500KB",
				directory:  "./",
			},
			mockURLs:   []string{"https://example.com/file3.txt", "https://example.com/file4.txt"},
			expectFail: true,
		},
		{
			name: "Empty URL file",
			args: args{
				filePath:   "./empty_urls.txt",
				outputFile: "output_empty.txt",
				limit:      "",
				directory:  "./",
			},
			mockURLs:   []string{},
			expectFail: true,
		},
	}

	// Create the test files before running tests
	var createdFiles []string
	for _, tt := range tests {
		if len(tt.mockURLs) > 0 {
			err := createMockFile(tt.args.filePath, tt.mockURLs)
			if err != nil {
				t.Fatalf("Failed to create mock file %s: %v", tt.args.filePath, err)
			}
			createdFiles = append(createdFiles, tt.args.filePath)
		}
	}

	defer func() {
		// Cleanup mock files after tests
		cleanupMockFiles(createdFiles)
	}()

	// Silence output during tests
	restore, _ := silenceOutput()
	defer restore()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure directory exists before running the test
			err := os.MkdirAll(tt.args.directory, 0o755)
			if err != nil {
				t.Fatalf("Failed to create directory: %v", err)
			}

			// Call the function
			DownloadMultipleFiles(tt.args.filePath, tt.args.outputFile, tt.args.limit, tt.args.directory)

			if !tt.expectFail && err == nil {
				// Check if the output file exists
				outputPath := filepath.Join(tt.args.directory, tt.args.outputFile)
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Expected output file at: %s, but it does not exist", outputPath)
				}
			}
		})
	}
}
