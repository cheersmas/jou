package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_DRIVER_NAME      = "sqlite3"
	DB_DATA_SOURCE_NAME = "lite.db"
)

func NewDatabase() (*sql.DB, error) {
	db, err := sql.Open(DB_DRIVER_NAME, DB_DATA_SOURCE_NAME)
	if err != nil {
		return nil, err
	}
	return db, nil
}
