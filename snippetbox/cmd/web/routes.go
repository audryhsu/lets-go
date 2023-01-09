package main

import "net/http"

// Update signature of routes() method so it returns a http.Handler instead of *http.ServeMux
// The routes() method returns a servemux containing the application router.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server to serve files out of "./ui/static".
	// Path given to http.Dir() is relative to project root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use mux.Handle() to register file server as handler for all URL
	// paths that start w/ "static". To match paths, strip "/static" prefix
	// before request reaches file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Wrap servemux with our middleware.
	// Last middleware func will call mux as the next middleware.
	// this is possible because our middleware takes a http.Handler. servemux implements the http.Handler interface because it has a ServeHttp() method.
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}