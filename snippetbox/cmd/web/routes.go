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

	// Non-auth routes use "dynamic" middleware chain plus CSRF check middleware
	// Create new middleware chain for middleware specifgic to dynamic app routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)

	// httprouter package provides method-based routing, clean URLs, and more robust pattern-matching.
	// alice ThenFunc() returns http.Handler (instead http.HandlerFunc), so switch to registering the route using router.Handler()
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// Authenticated routes use a "protected" middleware chain that includes requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreateForm))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	// Create middleware chain containing 'standard' middleware, which is used for every request our app receives
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return 'standard' middleware chain, followed by router
	return standard.Then(router)
}
