package repository

import (
	"crypto/sha512"
	"database/sql"
	"log"
	"uuid"

	_ "modernc.org/sqlite"
)

type Credential struct {
	Username string
	Password [sha512.Size]byte
}

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
			password BLOB NOT NULL CHECK (LENGTH(password) == 64)
		)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
            token BLOB CHECK (LENGTH(token) == 16)
		)`)
	if err != nil {
		log.Fatal(err)
	}
}

func NewCredential(c Credential) error {
	_, err := db.Exec(
		`INSERT INTO credentials (username, password) VALUES (?, ?)`,
		c.Username, c.Password[:],
	)
	return err
}

func GetByUsername(username string) (Credential, error) {
	var c Credential
	var password []byte
	err := db.QueryRow(
		`SELECT username, password FROM credentials WHERE username = ?`,
		username,
	).Scan(&c.Username, &password)
	copy(c.Password[:], password)
	return c, err
}

func NewSession(token uuid.UUID) error {
	_, err := db.Exec(
		`INSERT INTO sessions (token) VALUES (?)`, token[:],
	)
	return err
}

func SessionExists(token uuid.UUID) error {
	var t uuid.UUID
	err := db.QueryRow(
		`SELECT token FROM sessions WHERE token = ?`,
		token[:],
	).Scan(&t)
	return err
}

func DeleteSession(token uuid.UUID) error {
	_, err := db.Exec(
		`DELETE FROM sessions WHERE token = ?`, token[:],
	)
	return err
}
