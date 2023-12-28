package db

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func TestOpenDB(t *testing.T) {
	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
