package main

import (
	"os"
	"reflect"
	"testing"
)

func Test_validateURL(t *testing.T) {
	type args struct {
		link string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid URL with http",
			args:    args{link: "http://example.com"},
			wantErr: false,
		},
		{
			name:    "Valid URL with https",
			args:    args{link: "https://example.com"},
			wantErr: false,
		},
		{
			name:    "Valid URL with query params",
			args:    args{link: "https://example.com/path?query=value"},
			wantErr: false,
		},
		{
			name:    "Invalid URL without scheme",
			args:    args{link: "example.com"},
			wantErr: true,
		},
		{
			name:    "Invalid URL with spaces",
			args:    args{link: "http://example. com"},
			wantErr: true,
		},
		{
			name:    "Empty URL",
			args:    args{link: ""},
			wantErr: true,
		},
		{
			name:    "Valid URL with port number",
			args:    args{link: "https://example.com:8080"},
			wantErr: false,
		},
		{
			name:    "Valid localhost URL",
			args:    args{link: "http://localhost:3000"},
			wantErr: false,
		},
		// {
		// 	name:    "Invalid URL with incomplete domain",
		// 	args:    args{link: "http://.com"},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateURL(tt.args.link); (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Inputs
		wantErr bool
	}{
		{
			name: "Valid case with URL and file",
			args: []string{"program", "-O=output.txt", "http://example.com"},
			want: Inputs{
				file: "output.txt",
				url:  "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid case with mirroring and convert-links",
			args: []string{"program", "--mirror", "--convert-links", "http://example.com"},
			want: Inputs{
				mirroring:        true,
				convertLinksFlag: true,
				url:              "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid case with exclude flag",
			args: []string{"program", "--mirror", "-X=.jpg", "http://example.com"},
			want: Inputs{
				mirroring:   true,
				excludeFlag: ".jpg",
				url:         "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid case with rate limit and URL",
			args: []string{"program", "--rate-limit=100K", "http://example.com"},
			want: Inputs{
				rateLimit: "100K",
				url:       "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid case with source file and background download",
			args: []string{"program", "-i=input.txt", "-B", "http://example.com"},
			want: Inputs{
				sourcefile:       "input.txt",
				workInBackground: true,
				url:              "http://example.com",
			},
			wantErr: false,
		},
		{
			name: "Valid case with mirroring, reject, and URL",
			args: []string{"program", "--mirror", "--reject=.png", "http://example.com"},
			want: Inputs{
				mirroring:  true,
				rejectFlag: ".png",
				url:        "http://example.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			input := ParseArgs()
			if !reflect.DeepEqual(input, tt.want) {
				t.Errorf("ParseArgs() = %v, want %v", input, tt.want)
			}
		})
	}
}
