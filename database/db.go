package database

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

const (
	DB_DRIVER_NAME      = "sqlite"
	DB_DATA_SOURCE_NAME = "lite.db"
)

var (
	instance *sql.DB
	once     sync.Once
	initErr  error
)

// NewDatabase returns a singleton instance of the database connection
func NewDatabase() (*sql.DB, error) {
	once.Do(func() {
		instance, initErr = sql.Open(DB_DRIVER_NAME, DB_DATA_SOURCE_NAME)
		if initErr != nil {
			return
		}

		// Test the connection
		if err := instance.Ping(); err != nil {
			initErr = err
			return
		}
	})

	return instance, initErr
}

// GetInstance returns the existing database instance without trying to create a new one
func GetInstance() *sql.DB {
	return instance
}
