package main

import (
	"bytes"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"snippetbox.audryhsu.com/internal/models/mocks"
	"testing"
	"time"
)

// newTestApplication instantiates a new application struct with mocked errorLog and infoLog methods
func newTestApplication(t *testing.T) *application {
	// Create an instance of the template cache.
	templateCache, err := NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	// And a form decoder.
	formDecoder := form.NewDecoder()
	// And a session manager instance. Note that we use the same settings as // production, except that we *don't* set a Store for the session manager. // If no store is set, the SCS package will default to using a transient // in-memory store, which is ideal for testing purposes.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{}, // use mock
		users:          &mocks.UserModel{},    // use mock
		templateCache:  templateCache,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
	}
}

// define a custom testServer type which embeds a httptest.Server instance
type testServer struct {
	*httptest.Server
}

// newTestServer helper initializes and returns a new instance of custom test server type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// create a new test server and pass in a handler
	// starts a https server and listens on randomly-chosen port for the duration of the test.
	ts := httptest.NewTLSServer(h)

	// initialize a cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// add cookie jar to test server client. any response cookies will be stored and sent with subsequent requests when using this client.
	ts.Client().Jar = jar

	// disable redirect following by setting custom CheckRedirect function. This is called whenever a 3XXX response is received by client. Returning http.ErrUseLastResponse forces client to immediately return the received response
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

// get method on custom testServer type makes GET requests to a given URL path using the test server client and returns response status code, headers, and body
func (ts *testServer) get(t *testing.T, urlPath string) (statusCode int, headers http.Header, body string) {
	// network address that test server is listening on is in ts.URL field. Make GET requests against the test server using a configured test server client (http.Client).
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	bytebody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(bytebody)
	return res.StatusCode, res.Header, string(bytebody)
}