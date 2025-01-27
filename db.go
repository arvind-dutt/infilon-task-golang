package main

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/cetec")
	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Successfully connected to the database")
}
