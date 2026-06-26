package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	apiKey := "test-secret-key"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	tests := []struct {
		name       string
		path       string
		authHeader string
		wantCode   int
		wantBody   string
	}{
		{
			name:       "valid API key",
			path:       "/api/v1/items",
			authHeader: "Bearer " + apiKey,
			wantCode:   http.StatusOK,
			wantBody:   "success",
		},
		{
			name:       "missing Authorization header",
			path:       "/api/v1/items",
			authHeader: "",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "{\"message\":null,\"data\":null,\"error\":\"unauthorized\"}\n",
		},
		{
			name:       "invalid API key",
			path:       "/api/v1/items",
			authHeader: "Bearer wrong-key",
			wantCode:   http.StatusUnauthorized,
			wantBody:   "{\"message\":null,\"data\":null,\"error\":\"unauthorized\"}\n",
		},
		{
			name:       "root path bypasses auth",
			path:       "/",
			authHeader: "",
			wantCode:   http.StatusOK,
			wantBody:   "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authed := Auth(apiKey)(handler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			authed.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("Auth() status = %d, want %d", w.Code, tt.wantCode)
			}
			if w.Body.String() != tt.wantBody {
				t.Errorf("Auth() body = %q, want %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}
