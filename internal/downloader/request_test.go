package downloader

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpRequest(t *testing.T) {
	tests := []struct {
		name       string
		serverFunc func(w http.ResponseWriter, r *http.Request)
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Successful request",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello, World!"))
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Not Found",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantStatus: http.StatusNotFound,
			wantErr:    false,
		},
		{
			name: "Server Error",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantStatus: http.StatusInternalServerError,
			wantErr:    false,
		},
		{
			name: "Check headers",
			serverFunc: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("User-Agent") != "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.Header.Get("Accept") != "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.Header.Get("Accept-Language") != "en-US,en;q=0.5" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if r.Header.Get("Connection") != "keep-alive" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.serverFunc))
			defer server.Close()

			got, err := HttpRequest(server.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.wantStatus {
				t.Errorf("HttpRequest() status = %v, want %v", got.StatusCode, tt.wantStatus)
			}

			// Don't forget to close the response body
			if got != nil && got.Body != nil {
				got.Body.Close()
			}
		})
	}
}
