package mirror

import (
	"io/ioutil"
	"os"
	"testing"
)

// Mock the Mock function for testing
func MockgetLocalPath(url string) string {
	// Simulate path conversion for the purpose of testing
	return "local/path/to/" + url
}

func Test_convertCSSURLs(t *testing.T) {
	type args struct {
		cssContent string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Single URL with quotes",
			args: args{
				cssContent: `background: url("https://example.com/image.png");`,
			},
			want: `background: url('example.com/image.png');`,
		},
		{
			name: "Single URL with no quotes",
			args: args{
				cssContent: `background: url(https://example.com/image.png);`,
			},
			want: `background: url('example.com/image.png');`,
		},
		{
			name: "Multiple URLs",
			args: args{
				cssContent: `background: url("https://example.com/image1.png"), url('https://example.com/image2.png');`,
			},
			want: `background: url('example.com/image1.png'), url('example.com/image2.png');`,
		},
		{
			name: "No URLs",
			args: args{
				cssContent: `background: color(red);`,
			},
			want: `background: color(red);`,
		},
		{
			name: "URL with single quotes",
			args: args{
				cssContent: `background: url('example.com/image.png');`,
			},
			want: `background: url('example.com/image.png');`,
		},
		{
			name: "Escaped characters in URL",
			args: args{
				cssContent: `background: url("https://example.com/image%20with%20spaces.png");`,
			},
			want: `background: url('example.com/image with spaces.png');`,
		},
		{
			name: "URL without scheme",
			args: args{
				cssContent: `background: url("example.com/image.png");`,
			},
			want: `background: url('example.com/image.png');`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertCSSURLs(tt.args.cssContent); got != tt.want {
				t.Errorf("convertCSSURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocalPath(t *testing.T) {
	type args struct {
		originalURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid HTTP URL",
			args: args{
				originalURL: "http://example.com/path/to/file.txt",
			},
			want: "example.com/path/to/file.txt",
		},
		{
			name: "Valid HTTPS URL",
			args: args{
				originalURL: "https://example.com/path/to/file.txt",
			},
			want: "example.com/path/to/file.txt",
		},
		{
			name: "URL with no scheme",
			args: args{
				originalURL: "//example.com/path/to/file.txt",
			},
			want: "example.com/path/to/file.txt",
		},
		{
			name: "Local absolute path",
			args: args{
				originalURL: "/path/to/file.txt",
			},
			want: "path/to/file.txt",
		},
		{
			name: "Relative path",
			args: args{
				originalURL: "path/to/file.txt",
			},
			want: "path/to/file.txt",
		},
		{
			name: "Invalid URL",
			args: args{
				originalURL: "invalid-url",
			},
			want: "invalid-url",
		},
		{
			name: "Empty string",
			args: args{
				originalURL: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLocalPath(tt.args.originalURL); got != tt.want {
				t.Errorf("getLocalPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeHTTP(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HTTP URL",
			args: args{
				url: "http://example.com/path/to/resource",
			},
			want: "example.com/path/to/resource",
		},
		{
			name: "HTTPS URL",
			args: args{
				url: "https://example.com/path/to/resource",
			},
			want: "example.com/path/to/resource",
		},
		{
			name: "Base HTTP URL",
			args: args{
				url: "http://example.com/",
			},
			want: "example.com/index.html",
		},
		{
			name: "Base HTTPS URL",
			args: args{
				url: "https://example.com/",
			},
			want: "example.com/index.html",
		},
		{
			name: "Base URL without trailing slash",
			args: args{
				url: "http://example.com",
			},
			want: "example.com/index.html",
		},
		{
			name: "Invalid URL",
			args: args{
				url: "ftp://example.com/resource",
			},
			want: "ftp://example.com/resource",
		},
		{
			name: "Local URL without protocol",
			args: args{
				url: "example.com/path",
			},
			want: "example.com/path",
		},
		{
			name: "Empty URL",
			args: args{
				url: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeHTTP(tt.args.url); got != tt.want {
				t.Errorf("removeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFolder(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid directory path",
			args: args{
				path: "./test_directory",
			},
			want: true,
		},
		{
			name: "Valid file path",
			args: args{
				path: "./test_file.txt",
			},
			want: false,
		},
		{
			name: "Non-existent path",
			args: args{
				path: "./non_existent_path",
			},
			want: false,
		},
		{
			name: "Empty path",
			args: args{
				path: "",
			},
			want: false,
		},
		{
			name: "Symbolic link to directory",
			args: args{
				path: "./test_symlink",
			},
			want: true,
		},
	}

	// Silence terminal output
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// Set up for tests
	if err := os.Mkdir("test_directory", 0o755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll("test_directory") // Clean up the directory

	if err := os.WriteFile("test_file.txt", []byte("test content"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove("test_file.txt") // Clean up the file

	if err := os.Symlink("test_directory", "test_symlink"); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}
	defer os.Remove("test_symlink") // Clean up the symlink

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFolder(tt.args.path); got != tt.want {
				t.Errorf("IsFolder() = %v, want %v", got, tt.want)
			}
		})
	}

	// Clean up the pipe
	w.Close()
	ioutil.ReadAll(r) // Discard output
	r.Close()
}
