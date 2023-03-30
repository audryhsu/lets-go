package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"snippetbox.audryhsu.com/ui"
)

// Update signature of routes() method, so it returns a http.Handler instead of *http.ServeMux
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
	// convert ui.Files embedded filesystem and convert it to a http.FS type to satisfy the http.FileSystem interface and create  file server handler.
	fileServer := http.FileServer(http.FS(ui.Files))

	// Static files are now contained in "static" folder of ui.Files embedded filesystem, so we no longer need to strip the prefix from the request URL. Any requests that start with /static/ can be passed directly to file server. ("static/css/main.css")
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// add /ping route
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Non-auth routes use "dynamic" middleware chain plus CSRF check middleware
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)

	// httprouter package provides method-based routing, clean URLs, and more robust pattern-matching.
	// alice ThenFunc() returns http.Handler (instead http.HandlerFunc), so switch to registering the route using router.Handler()
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
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
