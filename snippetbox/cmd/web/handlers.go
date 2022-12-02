package main

import (
	"fmt"
	"html/template"
	// "log"
	"net/http"
	"strconv"
)

// Change signature of home handler as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// Initialize a slice containing path to two files, base file first.
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}
	// read template files into a template set. If error, log detailed error message.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// home handler is a method against application, it can access struct's fileds
		app.serverError(w, err)
		return
	}

	// Write the template content of "base" as response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	_, _ = fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	_, _ = w.Write([]byte("Create a new snippet..."))
}