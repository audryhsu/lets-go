package models

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
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

// Insert adds a new record to Users table
func (u *UserModel) Insert(name, email, password string) error {
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// store hashed password
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	result, err := u.DB.Exec(stmt, name, email, string(hashedPW))
	if err != nil {
		// check if error is mysql error, if so, check if matches duplicate error code
		// https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_dup_entry
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	// get the ID of newly inserted record
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

// Authenticate verifies whether user with email and password exists. Returns userID if valid.
func (u *UserModel) Authenticate(email, password string) (int, error) {
	var hashedPassword []byte
	var id int

	// if email doesn't exist in db, return error
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`
	row := u.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("couldn't find email in db")
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// if plaintext pw doesn't match hashed pw, return error
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		log.Println("badd password")
		return 0, ErrInvalidCredentials
	}

	log.Println("auth successful!")
	return id, nil
}

// Exists checks whether a user exists.
func (u *UserModel) Exists(id int) (bool, error) {
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	var exists bool
	err := u.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
