//go:build cgo

package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func connect(dbPath string) (*sql.DB, error) {
	// https://www.youtube.com/watch?v=XcAYkriuQ1o
	db, err := sql.Open("sqlite3", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, err
	}

	return db, nil
}
