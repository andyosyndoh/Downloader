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


