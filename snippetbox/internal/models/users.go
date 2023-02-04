package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// UserModel wraps a sql.DB connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to Useres table
func (u *UserModel) Insert(name, email, password string) (int, error) {
	stmt := `INSERT INTO users (name, email, hash, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// todo hash password with bcrypt
	hash := ""

	result, err := u.DB.Exec(stmt, name, email, hash)
	if err != nil {
		return 0, err
	}

	// get the ID of newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Authenticate verifies whether user with email and password exists. Returns userID if valid.
func (u *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks whether a user exists.
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
