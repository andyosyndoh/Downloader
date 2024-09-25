package downloader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOneDownload(t *testing.T) {
	// Test cases
	tests := []struct {
		name      string
		file      string
		url       string
		limit     string
		directory string
		wantFile  string
		wantError bool
	}{
		{
			name:      "Valid URL, default file name",
			file:      "meme.jpg",
			url:       "https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg",
			limit:     "",
			directory: "./test_download",
			wantFile:  "meme.jpg",
			wantError: false,
		},
		{
			name:      "No file name, derive from URL",
			file:      "",
			url:       "https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg",
			limit:     "",
			directory: "./test_download_no_file_name",
			wantFile:  "EMtmPFLWkAA8CIS.jpg",
			wantError: false,
		},
		{
			name:      "Rate limited download",
			file:      "meme_limited.jpg",
			url:       "https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg",
			limit:     "100k",
			directory: "./test_download_limited",
			wantFile:  "meme_limited.jpg",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the download
			OneDownload(tt.file, tt.url, tt.limit, tt.directory)

			// Check if the file was created in the correct directory
			downloadedFilePath := filepath.Join(tt.directory, tt.wantFile)
			_, err := os.Stat(downloadedFilePath)

			if tt.wantError {
				// Expecting an error, file should not exist
				if err == nil || !os.IsNotExist(err) {
					t.Errorf("Expected error but file was downloaded: %s", downloadedFilePath)
				}
			} else {
				// Expecting file to exist
				if os.IsNotExist(err) {
					t.Errorf("File not downloaded: %s", downloadedFilePath)
				} else {
					// Clean up after the test
					os.Remove(downloadedFilePath)
					os.RemoveAll(tt.directory)
				}
			}
		})
	}
}



func TestExpandPath(t *testing.T) {
	homeDir, _ := os.UserHomeDir() // Get the current user's home directory for testing

	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Expand home directory ~",
			args: args{path: "~/Documents"},
			want: filepath.Join(homeDir, "Documents"), // Expands ~/Documents to /home/user/Documents
		},
		{
			name: "Expand environment variable $HOME",
			args: args{path: "$HOME/Documents"},
			want: filepath.Join(homeDir, "Documents"), // Expands $HOME/Documents to /home/user/Documents
		},
		{
			name: "Expand relative path ./",
			args: args{path: "./"},
			want: func() string {
				absPath, _ := filepath.Abs(".")
				return absPath
			}(), // Converts ./ to absolute path
		},
		{
			name: "Expand relative path ../",
			args: args{path: "../"},
			want: func() string {
				absPath, _ := filepath.Abs("../")
				return absPath
			}(), // Converts ../ to absolute path
		},
		{
			name: "Absolute path remains unchanged",
			args: args{path: "/usr/local/bin"},
			want: "/usr/local/bin", // Absolute path stays the same
		},
		{
			name: "Empty path",
			args: args{path: ""},
			want: func() string {
				absPath, _ := filepath.Abs("")
				return absPath
			}(), // Converts empty string to current directory absolute path
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
