package models

import (
	"context"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"time"
)

type RegisterRequest struct {
	ID        string `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
}

func CreateUser(user RegisterRequest) (string, error) {
	db := database.GetDatabaseInstance()
	if db == nil || db.DB == nil {
		fmt.Println("Database connection error")
		log.Fatal("Database connection error")
		return "", fmt.Errorf("database connection error")
	}

	// Créer un contexte avec un délai d'expiration
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Requête SQL pour insérer un utilisateur
	query := `
        INSERT INTO users (id, nickname, first_name, last_name, email, age, gender, password)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// Préparer la déclaration SQL
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println("Failed to prepare create user statement: %v", err)
		return "", fmt.Errorf("failed to prepare create user statement: %v", err)
	}
	defer stmt.Close()

	// Exécuter la déclaration avec les données de l'utilisateur
	_, err = stmt.ExecContext(ctx, user.ID, user.Nickname, user.FirstName, user.LastName, user.Email, user.Age, user.Gender, user.Password)
	if err != nil {
		fmt.Println("Failed to create user: %v", err)
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	return user.ID, nil
}
