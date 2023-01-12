package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
	"time"
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

// methodNotAllowed helper is a convenience wrapper around clientError which sends 404 Not Found to user.
func (app *application) methodNotAllowed(w http.ResponseWriter) {
	app.clientError(w, http.StatusMethodNotAllowed)
}

// render method will retrieve appropriate template set from cache based on page (e.g. home.html)
// If no entry exists in cache with name, create a new error and call serverError()
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("template %s does not exist", page)
		app.serverError(w, err)
		return
	}
	// Write template to a buffer first to check for error.
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Safe to write to http.ResponseWriter if template is written to buffer w/o errors
	w.WriteHeader(status)

	// Write contents of buffer to http.ResponseWriter using WriteTo, which takes an io.Writer
	buf.WriteTo(w)
}

func (app *application) NewTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *application) decodePostForm(r *http.Request, dest any) error {
	// Parse request body to check it is well-formed, and if so, stores form data in r.PostForm map.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// call Decode() on decoder instance, passing the target destination as the first parameter
	if err := app.formDecoder.Decode(dest, r.PostForm); err != nil {
		// If we try to use an invalid target destination, Decode() method will return an error with the type *form.InvalidDecoderError. Use errors.As() to check for this specific error and panic instead of returning error.
		// why? if we pass something that isn't a non-nil pointer, this is a problem with our app code, not the user input, so we should handle this differently than returning 400.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// return err for all other types
		return err
	}
	return nil
}