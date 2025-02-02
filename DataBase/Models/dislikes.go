package models

import (
	"context"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"time"
)

func CreateDislike(dislike Liked_Post, user RegisterRequest) error {
	db := database.GetDatabaseInstance()
	if db == nil || db.DB == nil {
		log.Println("Database connection error")
		return fmt.Errorf("database connection error")
	}

	// Set a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if the dislike already exists
	var dislikeCount int
	err := db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID).Scan(&dislikeCount)

	if err != nil {
		log.Printf("Failed to check if dislike exists: %v\n", err)
		return err
	}

	// If dislike doesn't exist, check for a like
	if dislikeCount == 0 {
		// Check for existing like
		var likeCount int
		err = db.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID).Scan(&likeCount)
		if err != nil {
			log.Printf("Failed to check if like exists: %v\n", err)
			return err
		}

		// If a like exists, remove it
		if likeCount > 0 {
			_, err = db.DB.ExecContext(ctx, "DELETE FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID)
			if err != nil {
				log.Printf("Failed to delete like: %v\n", err)
				return err
			}
		}

		// Insert the dislike
		_, err = db.DB.ExecContext(ctx, "INSERT INTO disliked_posts (user_id, post_id) VALUES (?, ?)", user.ID, dislike.Post_ID)
		if err != nil {
			log.Printf("Failed to insert dislike: %v\n", err)
			return err
		}
	} else {
		// If the dislike exists, delete it
		_, err = db.DB.ExecContext(ctx, "DELETE FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID)
		if err != nil {
			log.Printf("Failed to delete dislike: %v\n", err)
			return err
		}
	}

	return nil
}

func CountDislikes(postID string) (int, error) {
	var count int
	db := database.GetDatabaseInstance()
	if db == nil || db.DB == nil {
		log.Println("Database connection error")
		return 0, fmt.Errorf("database connection error")
	}

	err := db.DB.QueryRow("SELECT COUNT(*) FROM disliked_posts WHERE post_id = ?", postID).Scan(&count)
	fmt.Println(count)
	return count, err
}
