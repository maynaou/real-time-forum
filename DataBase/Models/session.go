package models

import (
	"database/sql"
	"fmt"
	"time"

	database "handler/DataBase/Sqlite"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func CreateSession(session Session) (string, error) {
	query := "INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)"

	// Exécuter la requête directement
	_, err := database.DB.Exec(query, session.ID, session.UserID, time.Now(), session.ExpiresAt)
	if err != nil {
		fmt.Printf("Échec de la création de la session pour l'ID utilisateur : %s. Erreur : %v\n", session.UserID, err)
		return session.ID, err
	}

	return session.ID, nil
}

func GetSessionByID(id string) (Session, error) {
	var session Session
	query := "SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?"
	err := database.DB.QueryRow(query, id).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Session{}, nil
		}
		fmt.Printf("Failed to fetch session with ID: %s. Error: %v", id, err)
		return Session{}, err
	}
	return session, nil
}

func DeleteSession(id string) error {
	query := "DELETE FROM sessions WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	if err != nil {
		fmt.Printf("Échec de la suppression de la session avec l'ID : %s. Erreur : %v\n", id, err)
		return err
	}

	return nil
}
