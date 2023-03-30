package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"snippetbox.audryhsu.com/internal/assert"
	"testing"
)

func TestSnippetView(t *testing.T) {
	// test that GET requests to "/snippet/view/1" return 200 ok and req body contains expected content
	// GET requests to any other route return 404 not found
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	url := urlFormatter("/snippet/view")

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  url("1"),
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  url("509"),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  url("-1"),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  url("foo"),
			wantCode: http.StatusNotFound,
		},
		{
			name:     "empty id",
			urlPath:  url(""),
			wantCode: http.StatusNotFound,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			code, _, body := ts.get(t, c.urlPath)
			assert.Equal(t, code, c.wantCode)
			if c.wantBody != "" {
				mockSnippet, _ := app.snippets.Get(1)
				assert.StringContains(t, body, mockSnippet.Content)
			}
		})
	}
}

func urlFormatter(baseURL string) func(string) string {
	return func(param string) string {
		return fmt.Sprintf("%s/%s", baseURL, param)
	}
}

// TestUserSignup tests that a valid sign up form is submitted with a valid CSRF token.
func TestUserSignup(t *testing.T) {
	// send in post request with form in request body
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// send initial GET request to get sign up form with csrf token in request body
	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	// logs to test output if -v flag and test fails
	t.Logf("csrf token is: %q", validCSRFToken)

	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)
	// todo - empty username and email in use cases fail?
	tests := []struct {
		name        string
		wantCode    int
		username    string
		email       string
		password    string
		csrfToken   string
		wantFormTag string
	}{
		{
			name:      "valid signup",
			wantCode:  http.StatusSeeOther,
			username:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: validCSRFToken,
		},
		{
			name:      "no csrf token",
			wantCode:  http.StatusBadRequest,
			username:  validName,
			password:  validPassword,
			email:     validEmail,
			csrfToken: "wrongToken",
		},
		{
			name:        "empty form fields",
			wantCode:    http.StatusUnprocessableEntity,
			username:    "",
			email:       "",
			password:    "",
			wantFormTag: formTag,
			csrfToken:   validCSRFToken,
		},
		{
			name:        "empty form username",
			wantCode:    http.StatusUnprocessableEntity,
			username:    "",
			email:       validEmail,
			password:    validPassword,
			wantFormTag: formTag,
			csrfToken:   validCSRFToken,
		},
		{
			name:        "email not valid format",
			wantCode:    http.StatusUnprocessableEntity,
			username:    validName,
			email:       "alice12345",
			password:    validPassword,
			wantFormTag: formTag,
			csrfToken:   validCSRFToken,
		},
		{
			name:        "password is less than 8 characters",
			wantCode:    http.StatusUnprocessableEntity,
			username:    validName,
			email:       validEmail,
			password:    "123",
			csrfToken:   validCSRFToken,
			wantFormTag: formTag,
		},
		{
			name:        "email already in use",
			wantCode:    http.StatusUnprocessableEntity,
			username:    validName,
			email:       "dupe@example.com",
			password:    validPassword,
			wantFormTag: formTag,
			csrfToken:   validCSRFToken,
		},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", c.name)
			form.Add("email", c.email)
			form.Add("password", c.password)
			form.Add("csrf_token", c.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			assert.Equal(t, code, c.wantCode)

			if c.wantFormTag != "" {
				assert.StringContains(t, body, c.wantFormTag)
			}

		})
	}

}

func TestSnippetCreateForm(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// user must be authenticated via the requiresAuthentication middleware
	// if not authenticated, should redirect to /user/login page
	tests := []struct {
		name           string
		wantCode       int
		wantBody       string
		locationHeader string
	}{
		{
			name:           "unauthenticated",
			wantCode:       http.StatusSeeOther,
			locationHeader: "/user/login",
		},
		{
			name:     "authenticated",
			wantCode: http.StatusOK,
			wantBody: "<form action='/snippet/create' method='POST'>",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if test.name == "authenticated" {
				// mimic log in workflow by sending GET request to /user/login, extract CSRF token from response body
				// then make POST request to /user/login using mock user credentials
				_, _, body := ts.get(t, "/user/login")
				token := extractCSRFToken(t, body)

				// prints to test output if -v flag set and test fails
				t.Logf("token: %v", token)

				fmt.Printf("%v", body)
			}

			code, headers, _ := ts.get(t, "/snippet/create")
			assert.Equal(t, code, test.wantCode)
			assert.Equal(t, headers.Get("Location"), test.locationHeader)

		})

	}
}

// end to end testing that uses testutils package for set up
func TestPing(t *testing.T) {
	// create new instance of app struct for testing
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get(t, "/ping")
	// assert against returned http.Response instead of a http.ResponseRecorder
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, body, "OK")

}

// initial unit test
func TestPing_UnitTest(t *testing.T) {
	rr := httptest.NewRecorder()

	// initialize dummy http.Request (from a client)
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// pass in httptest.ResponseRecorder instead of http.ResponseWriter, and dummy request
	ping(rr, r)

	// Result() method gets the http.Response generated by the handler
	res := rr.Result()

	// Check that status code written by handler was 200
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// check response body written by handler equals "OK"
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")

}
