package models

import (
	"context"
	"fmt"
	"log"
	"time"

	database "handler/DataBase/Sqlite"
)

type RegisterRequest struct {
	ID        string    `json:"id"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Age       string    `json:"age"`
	Gender    string    `json:"gender"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"last_seen"`
	Online    bool      `json:"online"`
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
        INSERT INTO users (id, nickname, first_name, last_name, email, age, gender, password,created_at, last_seen)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?,?,?)`

	// Préparer la déclaration SQL
	stmt, err := db.DB.PrepareContext(ctx, query)
	if err != nil {
		fmt.Println("Failed to prepare create user statement: %v", err)
		return "", fmt.Errorf("failed to prepare create user statement: %v", err)
	}
	defer stmt.Close()

	// Exécuter la déclaration avec les données de l'utilisateur
	_, err = stmt.ExecContext(ctx, user.ID, user.Nickname, user.FirstName, user.LastName, user.Email, user.Age, user.Gender, user.Password, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		fmt.Println("Failed to create user: %v", err)
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	return user.ID, nil
}

func GetAllUsers(userID string, onlineMap map[string]bool) ([]RegisterRequest, error) {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return nil, fmt.Errorf("database connection error")
	}
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []RegisterRequest
	query := "SELECT * FROM users where id != ?"
	rows, err := database.DB.QueryContext(context, query, userID)
	if err != nil {
		fmt.Printf("failed to fetch users: %v", err)
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user RegisterRequest
		err := rows.Scan(&user.ID, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.Gender, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			fmt.Printf("failed to scan user row: %v", err)
			return nil, fmt.Errorf("failed to scan user row: %v", err)
		}
		user.Online = onlineMap[user.Nickname]
		users = append(users, user)
		fmt.Println(users)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("error while iterating through rows: %v", err)
		return nil, fmt.Errorf("error while iterating through rows: %v", err)
	}

	return users, nil
}
