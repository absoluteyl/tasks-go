package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestRouter(t *testing.T) {
	t.Run("ConnectDB", testConnectDB)

}

func testConnectDB(t *testing.T) {
	var dbDriver = "sqlite3"
	var dbPath = "../../db/tasks.db"

	db := connectDB(dbDriver, dbPath)
	defer db.Close()

	err := db.Ping()
	assert.NoError(t, err)
}
