package models

import (
	"context"
	"database/sql"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func CreateSession(session Session) (string, error) {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Println("Database connection error")
		log.Fatal("Database connection error")
		return session.ID, fmt.Errorf("database connection error")
	}
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)"
	statement, err := database.DB.PrepareContext(context, query)
	if err != nil {
		fmt.Println("Failed to prepare session statement for User ID: %s. Error: %v", session.UserID, err)
		return session.ID, fmt.Errorf("failed to prepare session statement: %v", err)
	}

	_, err = statement.ExecContext(context, session.ID, session.UserID, time.Now(), session.ExpiresAt)
	if err != nil {
		fmt.Println("Failed to create session for User ID: %s. Error: %v", session.UserID, err)
		return session.ID, err
	}

	return session.ID, nil
}

func GetSessionByID(id string) (Session, error) {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return Session{}, fmt.Errorf("database connection error")
	}
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

func UpdateSession(session Session) error {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return fmt.Errorf("database connection error")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE sessions SET user_id=?, expires_at=? WHERE id=?"
	stmt, err := database.DB.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("Failed to prepare update session statement for ID: %s. Error: %v", session.ID, err)
		return err
	}

	_, err = stmt.ExecContext(ctx, &session.UserID, &session.ExpiresAt, &session.ID)
	if err != nil {
		fmt.Printf("Failed to update session with ID: %s. Error: %v", session.ID, err)
		return err
	}

	return nil
}

func DeleteSession(id string) error {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return fmt.Errorf("database connection error")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM sessions WHERE id = ? "
	stmt, err := database.DB.PrepareContext(ctx, query)
	if err != nil {
		fmt.Printf("Failed to prepare delete session statement for ID: %s. Error: %v", id, err)
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		fmt.Printf("Failed to delete session with ID: %s. Error: %v", id, err)
		return err
	}

	return nil
}
