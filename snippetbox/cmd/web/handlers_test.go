package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox.audryhsu.com/internal/assert"
	"testing"
)

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
