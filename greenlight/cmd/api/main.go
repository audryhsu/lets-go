package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// config holds config settings for the application
type config struct {
	port int
	env  string
}
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var conf config

	flag.IntVar(&conf.port, "port", 4000, "API server port")
	flag.StringVar(&conf.env, "env", "development", "Environment (development|staging|production")

	flag.Parse()
	// intialize new logger which writes messages to stdout prefixed with current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := application{
		config: conf,
		logger: logger,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// declare a new HTTP server with timeout settings
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Starting %s server on %s", conf.env, server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)
}
