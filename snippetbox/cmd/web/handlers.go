package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
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
	data := app.NewTemplateData(r)
	data.Snippet = snippet

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
	log.Println("snippet create post handler!")

	if err := app.decodePostForm(r, &form); err != nil {
		log.Print("couldn't decode post form")
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(form.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(form.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1,7, or 365")

	if !form.Valid() {
		log.Println("failed form validation")
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	log.Println("trying to insert the snippet")
	if err != nil {
		log.Println("couldn't insert snippet into database")
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

type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// userSignup handles GET requests and renders HTML form to sign a user up.
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = UserSignupForm{}

	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form UserSignupForm
	// parse form data into UserSignup struct
	if err := app.decodePostForm(r, &form); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// validate data
	form.CheckField(form.MinChars(form.Password, 8), "password", "Password must be at least 8 characters long")
	form.CheckField(form.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(form.NotBlank(form.Name), "name", "name cannot be blank")
	form.CheckField(form.NotBlank(form.Email), "email", "email cannot be blank")
	form.CheckField(form.Matches(form.Email, validator.EmailRX), "email", "email must be valid address")

	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	// insert new user into db; if email already exists, re-render form with field error
	err := app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")
			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, http.StatusBadRequest, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Otherwise,add confirmation flash to session and redirect to login page
	app.sessionManager.Put(r.Context(), "flash", "User signed up successfully")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// userLoginForm represents and holds the form data
type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// userLogin displays the login page
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

// userLoginPost authenticates and logs in user
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	//err := app.formDecoder.Decode(&form, r.PostForm)
	// same as?
	err := app.decodePostForm(r, &form)
	if err != nil {
		log.Println("form decode error on user login")
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// validation checks -- email and password are provided and formats are correct
	form.CheckField(form.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(form.NotBlank(form.Password), "password", "This field field cannot be blank")
	form.CheckField(form.Matches(form.Email, validator.EmailRX), "email", "This field must be valid email")

	if !form.Valid() {
		log.Println("form failed validation")
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}
	// Good practice to generate a new session ID when auth state or priv levels change for a user (e.g. login/logout)
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// add ID of current user to session so they are 'logged in'
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

// userLogoutPost renews the session ID and removes userid from session store
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "Logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

// about displays the about page
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	app.render(w, http.StatusOK, "about.html", data)
}
