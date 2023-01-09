package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

// Update signature of routes() method so it returns a http.Handler instead of *http.ServeMux
// The routes() method returns a http.Handler containing the application router.
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Create a handler func which wraps our notFound() helper, then assign it as custom handler for 404 Not Found response. Ensures all 404 responses are standardized between notFound() calls and 404's from httprouter when no url pattern is matched.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.methodNotAllowed(w)
	})
	// Create a file server to serve files out of "./ui/static".
	// Path given to http.Dir() is relative to project root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Register file server as handler for all URL paths that start w/ "static". To match paths, strip "/static" prefix before request reaches file server.

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// httprouter package provides method-based routing, clean URLs, and more robust pattern-matching.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateForm)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create middleware chain containing 'standard' middleware, which is used for every request our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return 'standard' middleware chain, followed by router
	return standard.Then(router)
}