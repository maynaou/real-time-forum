package dataBase

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "./RealTimeForum.db")
	if err != nil {
		fmt.Println("Error in open database")
		os.Exit(1)
		return
	}

	ActivateForeingKey := `
		PRAGMA foreign_keys = ON;
	`
	_, err = Db.Exec(ActivateForeingKey)
	if err != nil {
		fmt.Println("Error creating users table")
		os.Exit(1)
	}

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = Db.Exec(createUsersTable)
	if err != nil {
		fmt.Println("Error creating users table")
		os.Exit(1)
	}
}
