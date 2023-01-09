package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"snippetbox.audryhsu.com/internal/models"
	"strings"
	"unicode/utf8"

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
	data := app.NewTemplateData(r)
	data.Snippets = snippets
	// render template passing in templateData of the latest snippets
	app.render(w, http.StatusOK, "home.html", data)
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// http-router stores named parameters in request context.
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
	data := app.NewTemplateData(r)
	data.Snippet = snippet
	// render an instance of templateData struct holding snippet data
	app.render(w, http.StatusOK, "view.html", data)
}

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Parse request body to check it is well-formed, and if so, stores form data in r.PostForm map.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	expires, err := strconv.Atoi(r.PostFormValue("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create an instance of snippetCreateForm containing values from form and empty map for validation errors
	form := snippetCreateForm{
		Title:       r.PostFormValue("title"),
		Content:     r.PostFormValue("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 { // RuneCountInString counts characters; len(s) counts bytes
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}
	permittedValue := expires == 1 || expires == 7 || expires == 365
	if !permittedValue {
		form.FieldErrors["expires"] = "This field must equal 1, 7, or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// use clean URL format in redirects
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// snippetCreateForm handles GET requests and renders HTML form to create snippets.
func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)

	// Initialize a new snippetCreateForm instance and pass to template
	// Without initializing the form field, the server will error out bc template cannot render nil as .Form in HTML
	data.Form = snippetCreateForm{Expires: 365}

	app.render(w, http.StatusOK, "create.html", data)
}