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
	t.Run("SetupRouter", testSetupRouter)
}

func testConnectDB(t *testing.T) {
	var dbDriver = "sqlite3"
	var dbPath = "../../db/tasks.db"

	db := connectDB(dbDriver, dbPath)
	defer db.Close()

	err := db.Ping()
	assert.NoError(t, err)
}

func testSetupRouter(t *testing.T) {
	var dbDriver = "sqlite3"
	var dbPath = "../../db/tasks.db"

	db := connectDB(dbDriver, dbPath)
	defer db.Close()

	mux := setupRouter(db)

	expectedRoutes := map[string][]string{
		"POST":   {"/auth", "/task"},
		"GET":    {"/tasks"},
		"PUT":    {"/task/:id"},
		"DELETE": {"/task/:id"},
	}

	for method, routes := range expectedRoutes {
		for i, route := range routes {
			assert.Equal(t, len(routes), len(mux.Routes[method]))
			assert.Equal(t, method, mux.Routes[method][i].Method)
			assert.Equal(t, route, mux.Routes[method][i].Path)
		}
	}

}
