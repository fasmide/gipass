package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

var Logins *DB

func init() {
	var err error
	Logins, err = NewLogins()
	if err != nil {
		log.Fatalf("unable to initate login db: %s", err)
	}
}

type DB struct {
	db *sql.DB
}

// NewLogins returns an opened database
func NewLogins() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to determine user home: %w", err)
	}

	p := path.Join(home, ".config/google-chrome", "Default", "Login Data")
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?immutable=1", p))
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	return &DB{db: db}, nil
}

func (d *DB) Query(s string) ([]Result, error) {
	stmt, err := d.db.Prepare("select origin_url, username_value, password_value, times_used from logins where origin_url LIKE ? ORDER BY times_used DESC")
	if err != nil {
		return nil, fmt.Errorf("unable to prepare query: %w", err)
	}

	r, err := stmt.Query(fmt.Sprintf("%%%s%%", s))
	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}

	results := make([]Result, 0)
	for r.Next() {
		p := Result{}
		err := r.Scan(&p.URL, &p.Username, &p.Password, &p.TimesUsed)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		results = append(results, p)
	}

	return results, nil
}
