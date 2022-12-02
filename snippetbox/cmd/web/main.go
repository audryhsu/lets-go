package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold app-wide dependencies.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// add 'addr' as cmd line flag, default value of ':4000
	// flag.Int(), flag.Bool()...auto convert flag value to type
	addr := flag.String("addr", ":4000", "HTTP network address")

	// parse cmd line flags and assign to addr variable.
	// must be called BEFORE using addr variable or else will contain default ':4000'
	flag.Parse()

	// create a new logger with custom prefix that writes to stdout
	// includes local date and time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// use log.Lshortfile to incl relevant file name and line number
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application struct with dependencies injected
	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	// initialize new http.Server struct and pass in our custom error logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Call app.routes() to get servemux containing our routes
	}
	infoLog.Printf("Starting server on %s", *addr)
	// call ListenAndServe() method on new http.Server struct
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}