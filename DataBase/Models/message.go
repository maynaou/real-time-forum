package models

import (
	"context"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"time"
)

type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
		IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`

}

func CreateMessage(message Message) (string, error) {
	database :=database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return message.ID, fmt.Errorf("database connection error")
	}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO messages (id,sender_id, receiver_id,content, is_read, create_at) VALUES (?, ?, ?, ?, ?, ?)"
	statement, err := database.DB.PrepareContext(context, query)
	if err != nil {
		fmt.Printf("Failed to prepare create message statement for sender: %s, receiver: %s. Error: %v", message.SenderID, message.ReceiverID, err)
		return message.ID, fmt.Errorf("failed to prepare create message statement: %v", err)
	}

	_, err = statement.ExecContext(context, &message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.IsRead, time.Now().UTC())
	if err != nil {
		fmt.Printf("Failed to create message for sender: %s, receiver: %s. Error: %v", message.SenderID, message.ReceiverID, err)
		return message.ID, fmt.Errorf("failed to create message: %v", err)
	}

	return message.ID, nil
}
