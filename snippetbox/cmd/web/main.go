package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Create a file server to serve files out of "./ui/static".
	// Path given to http.Dir() is relative to project root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use mux.Handle() to register file server as handler for all URL
	// paths that start w/ "static". To match paths, strip "/static" prefix
	// before request reaches file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}