package models

import (
	"fmt"

	database "handler/DataBase/Sqlite"
)

type LoginRequest struct {
	UserID   string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserDetails(req LoginRequest) (storedreq LoginRequest) {
	var identifier string
	if req.Nickname != "" {
		identifier = req.Nickname
	} else if req.Email != "" {
		identifier = req.Email
	} 

	query := "SELECT id, nickname, email, password FROM users WHERE nickname = ? OR email = ?"
	row := database.DB.QueryRow(query, identifier, identifier)

	err := row.Scan(&storedreq.UserID, &storedreq.Nickname, &storedreq.Email, &storedreq.Password)
	if err != nil {
		fmt.Printf("Failed to scan row for identifier %s: %v", identifier, err)
		return storedreq
	}

	return storedreq
}
