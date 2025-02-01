package models

import (
	"fmt"
	database "handler/DataBase/Sqlite"

	"log"
)

type LoginRequest struct {
	UserID   string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserDetails(req LoginRequest) (storedreq LoginRequest, err error) {
	db := database.GetDatabaseInstance()
	if db == nil || db.DB == nil {
		fmt.Println("Database connection error")
		log.Fatal("Database connection error")
		return
	}

	var identifier string
	if req.Nickname != "" {
		identifier = req.Nickname
	} else if req.Email != "" {
		identifier = req.Email
	} else {
		fmt.Println("Both username and email are empty")
		return storedreq, fmt.Errorf("both username and email are empty")
	}
	query := "SELECT id, nickname, email, password FROM users WHERE nickname = ? OR email = ?"
	row := db.DB.QueryRow(query, identifier, identifier)

	err = row.Scan(&storedreq.UserID, &storedreq.Nickname, &storedreq.Email, &storedreq.Password)
	if err != nil {

		fmt.Printf("Failed to scan row for identifier %s: %v", identifier, err)
		return storedreq, err
	}

	return storedreq, nil
}
