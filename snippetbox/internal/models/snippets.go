package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. Fields of struct correspond to the fields in MySQL snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

// Get Return a specific snippet based on id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// Latest Return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}