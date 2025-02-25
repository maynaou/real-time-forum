package models

import (
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
	Sender    bool      `json:"sender"`
}

func CreateUser(user RegisterRequest) (string, error) {
	query := `
        INSERT INTO users (id, nickname, first_name, last_name, email, age, gender, password, created_at, last_seen)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Exécuter la déclaration avec les données de l'utilisateur
	_, err := database.DB.Exec(query, user.ID, user.Nickname, user.FirstName, user.LastName, user.Email, user.Age, user.Gender, user.Password, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		log.Printf("Échec de la création de l'utilisateur : %v", err)
		return "", fmt.Errorf("échec de la création de l'utilisateur : %w", err)
	}

	return user.ID, nil
}

func GetAllUsers(onlineMap map[string]bool) ([]RegisterRequest, error) {
	var users []RegisterRequest
	query := "SELECT * FROM users"

	// Exécuter la requête pour obtenir les utilisateurs
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Échec de la récupération des utilisateurs : %v", err)
		return nil, fmt.Errorf("échec de la récupération des utilisateurs : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user RegisterRequest
		err := rows.Scan(&user.ID, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.Gender, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Printf("Échec de la lecture de la ligne de l'utilisateur : %v", err)
			return nil, fmt.Errorf("échec de la lecture de la ligne de l'utilisateur : %w", err)
		}
		user.Online = onlineMap[user.Nickname]
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erreur lors de l'itération à travers les lignes : %v", err)
		return nil, fmt.Errorf("erreur lors de l'itération à travers les lignes : %w", err)
	}

	return users, nil
}
