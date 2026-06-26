package middleware

import (
	"net/http"
	"strings"

	"github.com/spikycham/finance/network"
)

func Auth(apiKey string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				next.ServeHTTP(w, r)
				return
			}

			auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

			if auth != apiKey {
				network.ResponseError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
