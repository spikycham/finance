package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChain(t *testing.T) {
	var order []string

	m1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m1-before")
			next.ServeHTTP(w, r)
			order = append(order, "m1-after")
		})
	}

	m2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "m2-before")
			next.ServeHTTP(w, r)
			order = append(order, "m2-after")
		})
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler, m1, m2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	chained.ServeHTTP(w, req)

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("Chain() executed %d steps, want %d", len(order), len(expected))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("Chain() step %d = %q, want %q", i, order[i], v)
		}
	}
}

func TestChainNoMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	chained := Chain(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	chained.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Chain() with no middleware status = %d, want %d", w.Code, http.StatusOK)
	}
}
