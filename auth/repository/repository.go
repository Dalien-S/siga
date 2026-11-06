package repository

import (
	"database/sql"
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

type Credential struct {
	Username string
	Password string
}

var ErrNotFound = errors.New("not found")

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "db.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS credentials (
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)`)
	if err != nil {
		log.Fatal(err)
	}
}

func NewCredential(c Credential) error {
	_, err := db.Exec(
		`INSERT INTO credentials (username, password) VALUES (?, ?)`,
		c.Username, c.Password,
	)
	return err
}

func GetByUsername(username string) (Credential, error) {
	var c Credential
	err := db.QueryRow(
		`SELECT username, password FROM credentials WHERE username = ?`,
		username,
	).Scan(&c.Username, &c.Password)
	return c, err
}
