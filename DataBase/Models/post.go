package models

import (
	"context"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"strings"
	"time"
)

type Post struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Username   string    `json:"username"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Categories []string  `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
}

func CreatePost(post Post, user RegisterRequest) (string, error) {
	db := database.GetDatabaseInstance()
	if db == nil || db.DB == nil {
		fmt.Println("Database connection error")
		log.Fatal("Database connection error")
		return post.ID, fmt.Errorf("database connection error")
	}
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Che

	query := "INSERT INTO posts (id,user_id,title,content,category,created_at) VALUES (?, ?, ?, ?, ?, ?)"
	statement, err := db.DB.PrepareContext(context, query)
	if err != nil {
		fmt.Printf("Failed to prepare create post statement for user ID: %s, Title: %s, Error: %v", user.ID, post.Title, err)
		return post.ID, fmt.Errorf("failed to prepare create post statement: %v", err)
	}

	categoriesStr := strings.Join(post.Categories, ",")
	_, err = statement.ExecContext(context, &post.ID, &user.ID, &post.Title, &post.Content, categoriesStr, &post.CreatedAt)
	if err != nil {
		fmt.Printf("Failed to create post for user ID: %s, Title: %s. Error: %v", user.ID, post.Title, err)
		return post.ID, fmt.Errorf("failed to create post: %v", err)
	}

	return post.ID, nil

}

func GetAllPosts() ([]Post, error) {
	db := database.GetDatabaseInstance()

	if db == nil || db.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return nil, fmt.Errorf("database connection error")
	}
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT
		posts.id,
		posts.user_id,
		users.nickname,
		posts.title,
		posts.content,
		posts.category,
		posts.created_at,
		(SELECT COUNT(*) FROM liked_posts WHERE post_id = posts.id) AS likes,
		(SELECT COUNT(*) FROM disliked_posts WHERE post_id = posts.id) AS dislikes
	FROM
		posts
	INNER JOIN
		users
	ON
		posts.user_id = users.id
		 ORDER BY
        posts.created_at DESC
`
	rows, err := db.DB.QueryContext(context, query)
	if err != nil {
		fmt.Printf("Failed to execute query to get all posts ,Error: %v", err)
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	defer rows.Close()
	posts := make([]Post, 0)

	for rows.Next() {
		var post Post
		var categoriesStr string

		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &categoriesStr, &post.CreatedAt, &post.Likes, &post.Dislikes)
		if err != nil {
			fmt.Printf("Failed to scan row during post retrieval. Error: %v", err)
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		post.Categories = strings.Split(categoriesStr, ",")

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Failed to retrieve rows during post retrieval. Error: %v", err)
		return nil, fmt.Errorf("failed to retrieve rows: %v", err)
	}

	return posts, nil
}
