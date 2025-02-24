package models

import (
	"database/sql"
	"fmt"

	database "handler/DataBase/Sqlite"
)

type MessageData struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Message   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"`
}

func CreateMessage(username string, messageData MessageData) error {
	db := database.GetDatabaseInstance()
	if db.DB == nil {
		return fmt.Errorf("database connection error")
	}

	_, err := db.DB.Exec(`
		INSERT INTO messages (sender, receiver, content,created_at)
		VALUES (?, ?, ?,?)
	`, username, messageData.Receiver, messageData.Message, messageData.CreatedAt)

	return err
}

func GetMessages(sender, receiver, before string) ([]MessageData, error) {
	db := database.GetDatabaseInstance()
	if db.DB == nil {
		return nil, fmt.Errorf("database connection error")
	}
	fmt.Println(before, "LLLL")
	var rows *sql.Rows
	var err error

	if before != "" {
		rows, err = db.DB.Query(`
            SELECT sender, receiver, content, created_at
            FROM messages
            WHERE ((sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?))
            AND created_at < ?
            ORDER BY created_at DESC
            LIMIT 10
        `, sender, receiver, receiver, sender, before)
	} else {
		rows, err = db.DB.Query(`
            SELECT sender, receiver, content, created_at
            FROM messages
            WHERE (sender = ? AND receiver = ?) OR (sender = ? AND receiver = ?)
            ORDER BY created_at DESC
            LIMIT 10
        `, sender, receiver, receiver, sender)
	}

	if err != nil {
		fmt.Println("HHHH")
		return nil, err
	}

	defer rows.Close()

	var messages []MessageData
	for rows.Next() {
		var message MessageData
		if err := rows.Scan(&message.Sender, &message.Receiver, &message.Message, &message.CreatedAt); err != nil {
			fmt.Printf("Failed to scan row during message retrieval. Error: %v", err)
			continue
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("JJJJ")
		return nil, err
	}

	return messages, nil
}
