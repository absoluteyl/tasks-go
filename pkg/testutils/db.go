package testutils

import (
	"database/sql"
	"fmt"
	"os"
)

var DB_DRIVER = "sqlite3"
var DB_NAME = "test.db"

func ConnectDB() (*sql.DB, error) {
	var err error
	db, err := sql.Open(DB_DRIVER, DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to test database: %v", err)
	}

	return db, nil
}

func PrepareTaskTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status INTEGER DEFAULT 0
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("Error creating tasks table: %v", err)
	}

	return nil
}

func RemoveDB() error {
	err := os.Remove(DB_NAME)
	if err != nil {
		return fmt.Errorf("Error removing test database: %v", err)
	}

	return nil
}
