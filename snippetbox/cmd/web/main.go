package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox.audryhsu.com/internal/models"
	"time"
)

// Define an application struct to hold app-wide dependencies.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel // now model is available to handlers
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for MySQL DSN string
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL data source name")

	// parse cmd line flags and assign to addr variable.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Pass in the DSN from command line flag
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// close connection pool before main() function exits.
	defer db.Close()
	// initialize new template cache to add to app dependencies
	templateCache, err := NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// create an instance of form decoder and inject as app dependency
	formDecoder := form.NewDecoder()

	// initialize a new session manager and configure it to use our MySQL db as session store
	// set lifetime of 12 hours (automatic expiry)
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:        infoLog,
		errorLog:       errorLog,
		snippets:       &models.SnippetModel{DB: db}, // initialize a SnippetModel instance
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Call app.routes() to get servemux containing our routes
	}
	infoLog.Printf("Starting server on %s", *addr)
	// call ListenAndServe() method on new http.Server struct
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// create a connection and check for any errors
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}