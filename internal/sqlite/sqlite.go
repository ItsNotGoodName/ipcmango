package sqlite

import (
	"database/sql"
	"fmt"
)

func New(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("invalid database path: %s", dbPath)
	}

	db, err := connect(dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
