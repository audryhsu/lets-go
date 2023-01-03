package main

import (
	"errors"
	"fmt"
	"snippetbox.audryhsu.com/internal/models"

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

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}
	//// Initialize a slice containing path to two files, base file first.
	//files := []string{
	//	"./ui/html/base.html",
	//	"./ui/html/partials/nav.html",
	//	"./ui/html/pages/home.html",
	//}
	//// read template files into a template set. If error, log detailed error message.
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	// home handler is a method against application, it can access struct's fileds
	//	app.serverError(w, err)
	//	return
	//}
	//
	//// Write the template content of "base" as response body
	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// use SnippetModel object's Get method to retrieve snipped by ID. Return 404 not found if no matching record.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Write the snippet data as plain-text HTTP response body
	_, _ = fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title, content, expires := "0 snail", "0 snail\nClimb Mount Fuji,\nBut slowly,!\n\n -Kobayashi Issa", 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}