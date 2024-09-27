package downloader

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput captures both stdout and stderr output from a function.
func captureOutput(f func()) (string, string) {
	// Save the original stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create pipes to capture stdout and stderr
	stdoutReader, stdoutWriter, _ := os.Pipe()
	stderrReader, stderrWriter, _ := os.Pipe()

	// Redirect stdout and stderr to the pipes
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	// Buffer to read combined stdout and stderr output
	out := make(chan string)

	go func() {
		var bufStdout, bufStderr bytes.Buffer
		_, _ = ioutil.ReadAll(io.TeeReader(stdoutReader, &bufStdout))
		_, _ = ioutil.ReadAll(io.TeeReader(stderrReader, &bufStderr))
		out <- bufStdout.String() + bufStderr.String()
	}()

	// Run the function whose output is being captured
	f()

	// Close pipes and restore stdout and stderr
	stdoutWriter.Close()
	stderrWriter.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return <-out, "" // Return combined output
}

// Test cases for OneDownload// TestOneDownload tests the OneDownload function by capturing output.
func TestOneDownload(t *testing.T) {
	type args struct {
		file      string
		url       string
		limit     string
		directory string
	}
	tests := []struct {
		name       string
		args       args
		shouldFail bool
		expected   string
	}{
		{
			name: "Unsuccessful download with default directory and no limit",
			args: args{
				file:      "testfile.txt",
				url:       "https://example.com/testfile.txt",
				limit:     "",
				directory: "./downloads",
			},
			expected:   "Error: status 404 Not Found", // Expected output
			shouldFail: true,
		},
		{
			name: "Unsuccessful download with rate limit",
			args: args{
				file:      "testfile_rate.txt",
				url:       "https://example.com/testfile_rate.txt",
				limit:     "500KB",
				directory: "./downloads",
			},
			expected:   "Error: status 500 Internal Server Error", // Expected output
			shouldFail: true,
		},
		{
			name: "Download with missing file name (auto-generated)",
			args: args{
				file:      "",
				url:       "https://example.com/autogen.txt",
				limit:     "",
				directory: "./downloads",
			},
			expected:   "Error: status 500 Internal Server Error", // Expected output
			shouldFail: true,
		},
		{
			name: "Download with invalid URL",
			args: args{
				file:      "invalidfile.txt",
				url:       "https://invalid-url.com/file.txt",
				limit:     "",
				directory: "./downloads",
			},
			expected:   "Error downloading file: error sending request", // Expected output
			shouldFail: true,
		},
		{
			name: "Download to an invalid directory",
			args: args{
				file:      "testfile.txt",
				url:       "https://example.com/testfile.txt",
				limit:     "",
				directory: "/invalid/directory",
			},
			expected:   "Error: status 404 Not Found", // Expected output
			shouldFail: true,
		},
		{
			name: "Download with empty directory and file",
			args: args{
				file:      "",
				url:       "https://example.com/file.txt",
				limit:     "",
				directory: "",
			},
			expected:   "Error: status 404 Not Found", // Expected output
			shouldFail: true,
		},
		{
			name: "Download with valid URL and no limit",
			args: args{
				file:      "validfile.txt",
				url:       "https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg",
				limit:     "",
				directory: "./downloads",
			},
			expected:   "Download completed successfully", // Expected output
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture the output printed to the terminal (stdout + stderr)
			output, _ := captureOutput(func() {
				OneDownload(tt.args.file, tt.args.url, tt.args.limit, tt.args.directory)
			})

			// Check if the output contains the expected message
			if tt.shouldFail && !strings.Contains(output, tt.expected) {
				t.Errorf("Test %s failed, expected output: %s, got: %s", tt.name, tt.expected, output)
			}
			// For cases that should pass, ensure there's no error output
			if !tt.shouldFail && strings.Contains(output, "Error") {
				t.Errorf("Test %s failed, expected no error, got: %s", tt.name, output)
			}
		})
	}
	defer os.RemoveAll("./downloads")
}

func TestExpandPath(t *testing.T) {
	type args struct {
		path string
	}
	homeDir, _ := os.UserHomeDir() // Fetch the user's home directory for `~` expansion tests
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Expand home directory ~",
			args: args{
				path: "~/Documents/project",
			},
			want: filepath.Join(homeDir, "Documents/project"),
		},
		{
			name: "Expand environment variable $HOME",
			args: args{
				path: "$HOME/Documents/project",
			},
			want: filepath.Join(homeDir, "Documents/project"),
		},
		{
			name: "Expand environment variable $USER (should be in path)",
			args: args{
				path: "/home/$USER/project",
			},
			want: os.ExpandEnv("/home/$USER/project"), // Expect the $USER to be expanded correctly
		},
		{
			name: "Relative path ./",
			args: args{
				path: "./subfolder/file.txt",
			},
			want: func() string {
				abs, _ := filepath.Abs("./subfolder/file.txt")
				return abs
			}(),
		},
		{
			name: "Relative path ../",
			args: args{
				path: "../parentfolder/file.txt",
			},
			want: func() string {
				abs, _ := filepath.Abs("../parentfolder/file.txt")
				return abs
			}(),
		},
		{
			name: "Absolute path",
			args: args{
				path: "/usr/local/bin",
			},
			want: "/usr/local/bin", // Absolute paths should remain unchanged
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandPath(tt.args.path); got != tt.want {
				t.Errorf("ExpandPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
