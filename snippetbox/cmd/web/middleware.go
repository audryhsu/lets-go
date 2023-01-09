package main

import (
	"fmt"
	"net/http"
)

// secureHeaders sets Http security heads
func secureHeaders(next http.Handler) http.Handler {
	// http.HandlerFunc adapts a regular function into a http handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

// logRequest records the IP address of user and URL and method being requested. Method on app struct.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function which will ALWAYS be run in the event of a panic as Go unwinds the stack
		defer func() {
			// use builtin recover func to check if there has been a panic or not.
			// recover returns whatever the parameter passed to panic() was
			if err := recover(); err != nil {
				// set a Connection: close header on response
				w.Header().Set("Connection", "close")
				// Call app.serverError helper method to return 500 response, and pass in a new error object
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}