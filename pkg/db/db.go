package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "your_username"
	password = "your-password"
	dbname   = "your_database_name"
)

var db *sql.DB

func InitDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable ",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("Failed to open database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to ping database: %v", err)
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}

func GetDB() *sql.DB {
	return db
}
