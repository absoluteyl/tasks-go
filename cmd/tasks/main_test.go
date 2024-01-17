package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	setup()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func setup() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatalf("Error opening test database: %v", err)
	}
	defer db.Close()
}

func teardown() {
	// 删除测试用 SQLite3 数据库
	err := os.Remove("test.db")
	if err != nil {
		log.Printf("Error removing test database: %v", err)
	}
}
