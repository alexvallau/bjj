package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func connectDB() (*sql.DB, error) {
	dsn := "root:password@tcp(127.0.0.1:3306)/bjjNotes"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
