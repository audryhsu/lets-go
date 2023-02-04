package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"snippetbox.audryhsu.com/internal/models"
	"snippetbox.audryhsu.com/internal/validator"
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
	// retrieve the value for the key "flash" from session store. This deletes key and value. If no matching key, returns an empty string.
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.NewTemplateData(r)
	data.Snippet = snippet
	data.Flash = flash // pass flash message to template

	// render an instance of templateData struct holding snippet data
	app.render(w, http.StatusOK, "view.html", data)
}

type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // anonymous Validator type; "-" means ignore field during decoding
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(form.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1,7, or 365")

	if !form.Valid() {
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

	// use Put() method to add key/value pair to session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
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
