package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet type to hold the data for an individual snippet. Fields of struct correspond to the fields in MySQL snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModelInterface describes the methods that our SnippetModel struct has; created so that our application can expect a type that implements this interface, including our mock.SnippetModel!
type SnippetModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL statement to execute
	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Exec() method on connection pool to execute and return some basic info about what happened when statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get the ID of newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// cast int64 to int type
	return int(id), nil
}

// Get Return a specific snippet based on id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE id = ? AND expires > UTC_TIMESTAMP()`
	s := &Snippet{} // initialize a pointer to a new zeroed Snippet struct

	// use QueryRow method on connection pool to execute SQL statement. Returns a pointer to a sql.Row object which holds the result from db.
	row := m.DB.QueryRow(stmt, id)

	// row.Scan() copies query results into our zeroed Snippet instance, which should be POINTERS.
	// number of args must be exactly same as num of cols returned by SQL statement.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// if query returns no rows, row.Scan() returns a sql.ErrNoRows error.
		// use errors.Is() to check for specific error. If row not found, we return our own ErrNoRecord
		// Why? "encapsulate the model completely, so that our application isn't concerned with the underlying data store or reliant on datastore-specific errors for its behaviors".
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Latest Return 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	// Query() on the connection pool to exec. SQL statement. Returns sql.Rows resultset.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// if resultset is open, then underlying db conn remains open... don't use up all our conns!
	defer rows.Close()

	var snippets []*Snippet

	// use rows.Next iterate through rows in resultset nad prepare each row to be scanned.
	// if iteration of all rows is successful, resultset auto closes itself and db conn.
	for rows.Next() {
		s := &Snippet{}
		// rows.Scan() copies values from each field in row to new Snippet object.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// call rows.Err() to retrieve any error encountered during iteration. Important! Don't assume successful iteration over entire rulseset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
