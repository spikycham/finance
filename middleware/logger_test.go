package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logged := Logger(handler)

	tests := []struct {
		name       string
		method     string
		path       string
		remoteAddr string
		userAgent  string
		wantCode   int
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			path:       "/api/v1/items",
			remoteAddr: "192.168.1.1:12345",
			userAgent:  "Mozilla/5.0",
			wantCode:   http.StatusOK,
		},
		{
			name:       "POST request",
			method:     http.MethodPost,
			path:       "/api/v1/create",
			remoteAddr: "10.0.0.1:54321",
			userAgent:  "curl/7.0",
			wantCode:   http.StatusOK,
		},
		{
			name:       "empty user agent",
			method:     http.MethodGet,
			path:       "/",
			remoteAddr: "127.0.0.1:8080",
			userAgent:  "",
			wantCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req.RemoteAddr = tt.remoteAddr
			req.Header.Set("User-Agent", tt.userAgent)
			w := httptest.NewRecorder()

			logged.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Logger() status = %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestResponseWriterStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		writeFn  func(w http.ResponseWriter)
		wantCode int
	}{
		{
			name: "default 200",
			writeFn: func(w http.ResponseWriter) {
				// Don't call WriteHeader
			},
			wantCode: http.StatusOK,
		},
		{
			name: "explicit 201",
			writeFn: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusCreated)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "explicit 404",
			writeFn: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "explicit 500",
			writeFn: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			rw := newResponseWriter(w)

			tt.writeFn(rw)

			if rw.statusCode != tt.wantCode {
				t.Errorf("statusCode = %d, want %d", rw.statusCode, tt.wantCode)
			}
		})
	}
}
