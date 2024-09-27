package background

import (
	"os"
	"testing"
)

func TestLoadShowProgressState(t *testing.T) {
	tests := []struct {
		name       string
		fileExists bool
		fileData   string
		want       bool
		wantErr    bool
		setup      func() // Optional setup function for each test
	}{
		{
			name:       "File does not exist",
			fileExists: false,
			want:       true, // Default value when file doesn't exist
			wantErr:    false,
		},
		{
			name:       "File exists with valid true value",
			fileExists: true,
			fileData:   "true",
			want:       true,
			wantErr:    false,
		},
		{
			name:       "File exists with valid false value",
			fileExists: true,
			fileData:   "false",
			want:       false,
			wantErr:    false,
		},
		{
			name:       "File exists with invalid value",
			fileExists: true,
			fileData:   "invalid_boolean",
			want:       false,
			wantErr:    true, // Expect an error due to invalid boolean parsing
		},
		{
			name:       "Error reading file",
			fileExists: true,
			fileData:   "", // Simulate error by leaving the file empty
			want:       false,
			wantErr:    true,
			setup: func() {
				// Optionally, you can mock or simulate an error reading the file here.
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Optionally run setup if needed
			if tt.setup != nil {
				tt.setup()
			}

			// Create or delete the test file based on the case
			if tt.fileExists {
				err := os.WriteFile(tempConfigFile, []byte(tt.fileData), 0o644)
				if err != nil {
					t.Fatalf("Failed to create temp config file: %v", err)
				}
			} else {
				_ = os.Remove(tempConfigFile) // Ensure the file does not exist
			}

			got, err := LoadShowProgressState()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadShowProgressState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadShowProgressState() = %v, want %v", got, tt.want)
			}

			// Clean up after each test
			_ = os.Remove(tempConfigFile)
		})
	}
}

// TestSaveShowProgressState tests the SaveShowProgressState function.
func TestSaveShowProgressState(t *testing.T) {
	type args struct {
		showProgress bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Save showProgress true",
			args:    args{showProgress: true},
			wantErr: false,
		},
		{
			name:    "Save showProgress false",
			args:    args{showProgress: false},
			wantErr: false,
		},
		{
			name:    "Fail to save due to permission error",
			args:    args{showProgress: true},
			wantErr: true,
		},
	}

	// Save the original file permissions for cleanup
	originalPerm, err := os.Stat(tempConfigFile)
	if err != nil {
		originalPerm = nil
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In case of a permission error test, change permissions to simulate the error
			if tt.name == "Fail to save due to permission error" {
				os.Chmod(tempConfigFile, 0o000) // Make the file unreadable/writable
			} else {
				// Ensure the file is writable
				if originalPerm != nil {
					os.Chmod(tempConfigFile, 0o644) // Restore permissions to writable
				}
			}

			if err := SaveShowProgressState(tt.args.showProgress); (err != nil) != tt.wantErr {
				t.Errorf("SaveShowProgressState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Cleanup: Restore original permissions and remove the test file
	if originalPerm != nil {
		os.Chmod(tempConfigFile, originalPerm.Mode().Perm())
	}
	os.Remove(tempConfigFile)
}

func TestDownloadInBackground(t *testing.T) {
	type args struct {
		file      string
		urlStr    string
		rateLimit string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Valid URL with output file",
			args: args{
				file:      "output.txt",
				urlStr:    "https://example.com/file.txt",
				rateLimit: "200k",
			},
		},
		{
			name: "Valid URL with default filename",
			args: args{
				file:      "",
				urlStr:    "https://example.com/image.png",
				rateLimit: "1M",
			},
		},
		{
			name: "Valid URL with empty rate limit",
			args: args{
				file:      "output.txt",
				urlStr:    "https://example.com/file.txt",
				rateLimit: "",
			},
		},
	}

	// Temporary log file for testing
	tempLogFile := "wget-log"
	// Cleanup log file before tests
	os.Remove(tempLogFile)

	// Redirect stdout and stderr to io.Discard
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	_, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DownloadInBackground(tt.args.file, tt.args.urlStr, tt.args.rateLimit)

			// For validation purposes, you can check the log file or the output file created
			if tt.name == "Invalid URL" {
				// Check that log file was created
				if _, err := os.Stat(tempLogFile); os.IsNotExist(err) {
					t.Log("Log file not created for invalid URL, as expected.")
				} else {
					t.Errorf("Log file should not be created for an invalid URL")
				}
			} else {
				// Check that log file was created for valid cases
				if _, err := os.Stat(tempLogFile); os.IsNotExist(err) {
					t.Errorf("Log file not created for valid URL")
				}
			}
		})
	}

	// Cleanup
	os.Remove(tempLogFile)
	os.Remove("progress_config.txt")

	// Restore original stdout and stderr
	w.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
}
