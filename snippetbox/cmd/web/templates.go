package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox.audryhsu.com/internal/models"
	"snippetbox.audryhsu.com/ui"
	"time"
)

// templateData is a holding structure for any dynamic data we want to pass to HTML templates.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// humanDate returns a nicely formatted string of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize template.FuncMap object and store it in a global variable. This is a lookup between names of custom template funcs and funcs themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// NewTemplateCache creates a cache of parsed templates ready for use by handler functions to render dynamic data. Each page (key) has a corresponding set of templates (value).
func NewTemplateCache() (map[string]*template.Template, error) {
	// initialize new map
	cache := map[string]*template.Template{}

	// Use fs.Glob to get slice of all filepaths in ui.Files embedded fs that match the pattern "./ui/html/pages/*.html" (e.g. all of the "page" templates)
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// extract file name (home.html) from full filepath of page
		name := filepath.Base(page)

		// create slice containing filepath patterns for templates we want to parse.
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		// use ParseFS() instead of ParseFiles() to parse the template files from ui.Files embedded fs
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// add template set to map, with page as key
		cache[name] = ts
	}
	return cache, nil
}
