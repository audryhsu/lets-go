package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"snippetbox.audryhsu.com/internal/models"

	// "log"
	"net/http"
	"strconv"
)

// Change signature of home handler as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.NewTemplateData()
	data.Snippets = snippets
	// render template passing in templateData of latest snippets
	app.render(w, http.StatusOK, "home.html", data)
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// httprouter stores named parameters in request context.
	params := httprouter.ParamsFromContext(r.Context())

	// use ByName() method to get value of "id" named param from slice and validate
	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.NewTemplateData()
	data.Snippet = snippet
	// render an instance of templateData struct holding snippet data
	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title, content, expires := "0 snail", "0 snail\nClimb Mount Fuji,\nBut slowly,!\n\n -Kobayashi Issa", 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// use clean URL format in redirects
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/:id", id), http.StatusSeeOther)
}

// snippetCreate handles GET requests and returns form to create snippets.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dsiplay form for creating a new snippet..."))
}