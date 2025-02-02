package models

import (
	"context"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"time"
)

type Liked_Post struct {
	Post_ID string `json:"post_id"`
	User_ID string `json:"user_id"`
}

func CreateLike(like Liked_Post, user RegisterRequest) error {
	db := database.GetDatabaseInstance()
	if db == nil || db.DB == nil {
		log.Println("Database connection error")
		return fmt.Errorf("database connection error")
	}

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the like already exists
	var likeCount int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID).Scan(&likeCount)

	if err != nil {
		log.Printf("Failed to check if like exists: %v\n", err)
		return err
	}

	// If like doesn't exist, check for a dislike
	if likeCount == 0 {
		// Check for existing dislike
		var dislikeCount int
		err = db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID).Scan(&dislikeCount)
		if err != nil {
			log.Printf("Failed to check if dislike exists: %v\n", err)
			return err
		}

		// If a dislike exists, remove it
		if dislikeCount > 0 {
			_, err = db.DB.ExecContext(ctx, "DELETE FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID)
			if err != nil {
				log.Printf("Failed to delete dislike: %v\n", err)
				return err
			}
		}

		// Insert the like
		_, err = db.DB.ExecContext(ctx, "INSERT INTO liked_posts (user_id, post_id) VALUES (?, ?)", user.ID, like.Post_ID)
		if err != nil {
			log.Printf("Failed to insert like: %v\n", err)
			return err
		}
	} else {
		// If the like exists, delete it
		_, err = db.DB.ExecContext(ctx, "DELETE FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID)
		if err != nil {
			log.Printf("Failed to delete like: %v\n", err)
			return err
		}
	}

	return nil
}

