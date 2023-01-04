package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.audryhsu.com/internal/models"
	"time"
)

// templateData is a holding structure for any dynamic data we want to pass to HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// humanDate returns a nicely formatted string of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2023 at 15:04")
}

// Initialize template.FuncMap object and store it in a global variable. This is a lookup between names of custom template funcs and funcs themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// NewTemplateCache creates a cache of parsed templates ready for use by handler functions to render dynamic data. Each page (key) has a corresponding set of templates (value).
func NewTemplateCache() (map[string]*template.Template, error) {
	// initialize new map
	cache := map[string]*template.Template{}

	// Get slice of all filepaths that match pattern "./ui/html/pages/*.html"
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// extract file name (home.html) from full filepath of page
		name := filepath.Base(page)

		// Create a new blank template, then register template.FuncMap before parsing base template into template set
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		// call ParseGlob() on THIS template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		// add page to this template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to map, with page as key
		cache[name] = ts
	}
	return cache, nil
}