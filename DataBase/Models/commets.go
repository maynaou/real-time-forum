package models

import (
	"context"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Author    string    `json:"author"`
	PostID    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func CreateComment(comment Comment, user RegisterRequest) (string, error) {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return comment.ID, fmt.Errorf("database connection error")
	}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO comments (id, user_id, post_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	statement, err := database.DB.PrepareContext(context, query)
	if err != nil {
		fmt.Printf("Failed to prepare create comment statement: %v", err)
		return comment.ID, fmt.Errorf("failed to prepare create comment statement: %v", err)
	}

	_, err = statement.ExecContext(context, comment.ID, user.ID, comment.PostID, comment.Content, time.Now().UTC(), time.Now().UTC())
	if err != nil {
		fmt.Printf("Failed to create comment: %v", err)
		return comment.ID, fmt.Errorf("failed to create comment: %v", err)
	}

	return comment.ID, nil
}

func GetCommentsByPostID(postID string) ([]Comment, error) {
	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return nil, fmt.Errorf("database connection error")
	}
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT comments.id, comments.user_id, users.username, comments.post_id, comments.content, comments.created_at, comments.updated_at " +
		"FROM comments " +
		"JOIN users ON comments.user_id = users.id " +
		"WHERE comments.post_id = ?"

	rows, err := database.DB.QueryContext(context, query, postID)
	if err != nil {
		fmt.Printf("Failed to execute query: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.Author, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			fmt.Printf("Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Failed to retrieve rows: %v", err)
		return nil, fmt.Errorf("failed to retrieve rows: %v", err)
	}

	return comments, nil
}
