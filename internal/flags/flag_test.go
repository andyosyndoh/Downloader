package flags

import (
	"io/ioutil"
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
			name:    "Valid HTTP URL",
			args:    args{link: "http://example.com"},
			wantErr: false,
		},
		{
			name:    "Valid HTTPS URL",
			args:    args{link: "https://example.com"},
			wantErr: false,
		},
		{
			name:    "Valid URL with path",
			args:    args{link: "https://example.com/path/to/resource"},
			wantErr: false,
		},
		{
			name:    "Valid URL with query parameters",
			args:    args{link: "https://example.com/search?q=test&page=1"},
			wantErr: false,
		},
		{
			name:    "Invalid URL - missing scheme",
			args:    args{link: "example.com"},
			wantErr: true,
		},
		{
			name:    "Invalid URL - empty string",
			args:    args{link: ""},
			wantErr: true,
		},
		{
			name:    "Invalid URL - malformed",
			args:    args{link: "http://[::1]:namedport"},
			wantErr: true,
		},
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
		name string
		args []string
		want Inputs
	}{
		{
			name: "Basic URL",
			args: []string{"program", "https://example.com"},
			want: Inputs{URL: "https://example.com"},
		},
		{
			name: "URL with output file",
			args: []string{"program", "-O=output.html", "https://example.com"},
			want: Inputs{URL: "https://example.com", File: "output.html"},
		},
		{
			name: "URL with path",
			args: []string{"program", "-P=/downloads", "https://example.com"},
			want: Inputs{URL: "https://example.com", Path: "/downloads"},
		},
		{
			name: "URL with rate limit",
			args: []string{"program", "--rate-limit=100k", "https://example.com"},
			want: Inputs{URL: "https://example.com", RateLimit: "100k"},
		},
		{
			name: "Mirror mode",
			args: []string{"program", "--mirror", "https://example.com"},
			want: Inputs{URL: "https://example.com", Mirroring: true},
		},
		{
			name: "Mirror mode with convert links",
			args: []string{"program", "--mirror", "--convert-links", "https://example.com"},
			want: Inputs{URL: "https://example.com", Mirroring: true, ConvertLinksFlag: true},
		},
		{
			name: "Mirror mode with reject",
			args: []string{"program", "--mirror", "--reject=.pdf", "https://example.com"},
			want: Inputs{URL: "https://example.com", Mirroring: true, RejectFlag: ".pdf"},
		},
		{
			name: "Mirror mode with exclude",
			args: []string{"program", "--mirror", "--exclude=private", "https://example.com"},
			want: Inputs{URL: "https://example.com", Mirroring: true, ExcludeFlag: "private"},
		},
		{
			name: "Background download",
			args: []string{"program", "-B", "https://example.com"},
			want: Inputs{URL: "https://example.com", WorkInBackground: true},
		},
		{
			name: "Input from file",
			args: []string{"program", "-i=urls.txt"},
			want: Inputs{Sourcefile: "urls.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			_, _ = captureOutput(func() {
				defer func() {
					if r := recover(); r != nil {
						if tt.want == (Inputs{}) {
							// This was expected, do nothing
						} else {
							t.Errorf("ParseArgs() panicked unexpectedly: %v", r)
						}
					}
				}()

				got := ParseArgs()
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ParseArgs() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

func captureOutput(f func()) (string, string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return string(out), ""
}
