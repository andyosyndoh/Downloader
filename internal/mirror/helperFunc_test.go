package mirror

import (
	"testing"
)

func Test_isRejected(t *testing.T) {
	type args struct {
		url         string
		rejectTypes string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "No reject types",
			args: args{url: "http://example.com/image.jpg", rejectTypes: ""},
			want: false,
		},
		{
			name: "Rejected file type",
			args: args{url: "http://example.com/document.pdf", rejectTypes: ".pdf,.doc,.txt"},
			want: true,
		},
		{
			name: "Not rejected file type",
			args: args{url: "http://example.com/image.jpg", rejectTypes: ".pdf,.doc,.txt"},
			want: false,
		},
		{
			name: "Multiple reject types, one match",
			args: args{url: "http://example.com/document.doc", rejectTypes: ".pdf,.doc,.txt"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRejected(tt.args.url, tt.args.rejectTypes); got != tt.want {
				t.Errorf("isRejected() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRejectedPath(t *testing.T) {
	type args struct {
		url         string
		pathRejects string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "No path rejects",
			args: args{url: "http://example.com/page", pathRejects: ""},
			want: false,
		},
		{
			name: "Rejected path",
			args: args{url: "http://example.com/admin/dashboard", pathRejects: "/admin,/private"},
			want: true,
		},
		{
			name: "Not rejected path",
			args: args{url: "http://example.com/public/page", pathRejects: "/admin,/private"},
			want: false,
		},
		{
			name: "Rejected path with multiple options",
			args: args{url: "http://example.com/private/profile", pathRejects: "/admin,/private,/settings"},
			want: true,
		},
		{
			name: "Invalid reject path (no leading slash)",
			args: args{url: "http://example.com/admin/dashboard", pathRejects: "admin,private"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRejectedPath(tt.args.url, tt.args.pathRejects); got != tt.want {
				t.Errorf("isRejectedPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		str    string
		substr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Substring present",
			args: args{str: "Hello, world!", substr: "world"},
			want: true,
		},
		{
			name: "Substring not present",
			args: args{str: "Hello, world!", substr: "universe"},
			want: false,
		},
		{
			name: "Empty substring",
			args: args{str: "Hello, world!", substr: ""},
			want: true,
		},
		{
			name: "Substring longer than string",
			args: args{str: "Short", substr: "This is longer"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.str, tt.args.substr); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fileExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Existing file",
			args: args{path: "./internal/flags/flag.go"},
			want: true,
		},
		{
			name: "Non-existing file",
			args: args{path: "testdata/non_existing_file.txt"},
			want: false,
		},
		{
			name: "Existing directory",
			args: args{path: "./internal"},
			want: true,
		},
		{
			name: "Non-existing directory",
			args: args{path: "testdata/non_existing_dir"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileExists(tt.args.path); got != tt.want {
				t.Errorf("fileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractDomain(t *testing.T) {
	type args struct {
		urlStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Valid URL with subdomain",
			args:    args{urlStr: "https://www.example.com/page"},
			want:    "www.example.com",
			wantErr: false,
		},
		{
			name:    "Valid URL without subdomain",
			args:    args{urlStr: "http://example.com"},
			want:    "example.com",
			wantErr: false,
		},
		{
			name:    "Valid URL with port",
			args:    args{urlStr: "https://example.com:8080/page"},
			want:    "example.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractDomain(tt.args.urlStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidAttribute(t *testing.T) {
	type args struct {
		tagName string
		attrKey string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid a href",
			args: args{tagName: "a", attrKey: "href"},
			want: true,
		},
		{
			name: "Valid img src",
			args: args{tagName: "img", attrKey: "src"},
			want: true,
		},
		{
			name: "Valid script src",
			args: args{tagName: "script", attrKey: "src"},
			want: true,
		},
		{
			name: "Valid link href",
			args: args{tagName: "link", attrKey: "href"},
			want: true,
		},
		{
			name: "Invalid tag",
			args: args{tagName: "div", attrKey: "class"},
			want: false,
		},
		{
			name: "Invalid attribute for valid tag",
			args: args{tagName: "a", attrKey: "class"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidAttribute(tt.args.tagName, tt.args.attrKey); got != tt.want {
				t.Errorf("isValidAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}


