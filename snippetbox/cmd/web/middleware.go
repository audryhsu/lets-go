package main

import (
	"context"
	"fmt"
	"github.com/justinas/nosurf"
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if not authenticated, redirect
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// otherwise set header so that pages require auth are not stored in users browser cache
		w.Header().Add("Cache-Control", "no-store")
		// call next handler in chain
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get authenticatedUserID from session data
		userId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

		// if don't have authenticated user, pass original request to next handler
		if userId == 0 { // check for nil value
			next.ServeHTTP(w, r)
			return
		}

		// if there is an auth user ID in session data, check db to see if user id exists in database
		if exists, err := app.users.Exists(userId); err != nil {
			app.serverError(w, err)
			return
		} else if exists {
			// update request context to include new context key indicated auth is good
			// create a copy of the request with new context
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// noSurf middleware function creates a CSRF handler that will call the next middleware if the check passes. Check uses a customized CSRF http cookie
func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
