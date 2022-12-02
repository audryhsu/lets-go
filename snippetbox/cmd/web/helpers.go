package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Logs the stack trace using errorLog and responds with a 500 Internal Server error
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	// replies to request with HTTP code and message
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// helper sends a specific status code and corresponding description to the user.
// ex: 400 "Bad Request" when there's a problem with user request.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound helper is a convenience wrapper around clientError which sends 404 Not Found to user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}