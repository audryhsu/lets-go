package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox.audryhsu.com/internal/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// create and pass a mock HTTP handler to secureHeaders middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("next handler was called"))
	})
	// Mock middleware chain
	// Because secureHeaders returns a http.Handler, call the ServeHTTP() method with response recorder and dummy req
	secureHeaders(next).ServeHTTP(rr, req)
	res := rr.Result()

	// Check that middleware correctly set headers on responses
	tests := []struct {
		header   string
		expected string
	}{
		{
			header:   "Content-Security-Policy",
			expected: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
		},
		{
			header:   "Referer-Policy",
			expected: "origin-when-cross-origin",
		},
		{
			header:   "X-Content-Type-Options",
			expected: "nosniff",
		},
		{
			header:   "X-Frame-Options",
			expected: "deny",
		},
		{
			header:   "X-XSS-Protection",
			expected: "0",
		},
	}
	for _, test := range tests {
		t.Run(test.header, func(t *testing.T) {
			assert.Equal(t, res.Header.Get(test.header), test.expected)
		})
	}
	// Check that middleware correctly called the next handler
	assert.Equal(t, res.StatusCode, http.StatusOK)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "next handler was called")

}
