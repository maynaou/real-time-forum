package models

import (
	"fmt"
	"html"
	"strings"
	"time"

	database "handler/DataBase/Sqlite"
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
	Comments   int       `json:"comments"`
	IsLiked    bool      `json:"isLiked"`
}

func CreatePost(post Post, user RegisterRequest) (string, error) {
	query := "INSERT INTO posts (id, user_id, title, content, category, created_at) VALUES (?, ?, ?, ?, ?, ?)"
	post.Content = html.EscapeString(post.Content)
	post.Title = html.EscapeString(post.Title)
	categoriesStr := strings.Join(post.Categories, ",")
	_, err := database.DB.Exec(query, post.ID, user.ID, post.Title, post.Content, categoriesStr, post.CreatedAt)
	if err != nil {
		fmt.Printf("Échec de la création du post pour l'ID utilisateur : %s, Titre : %s. Erreur : %v\n", user.ID, post.Title, err)
		return post.ID, fmt.Errorf("échec de la création du post : %v", err)
	}

	return post.ID, nil
}

func GetAllPosts(userID string) ([]Post, error) {
	likedPostIDs, err := GetLikedPostIDs(userID)
	if err != nil {
		return nil, err
	}

	likedPostMap := make(map[string]bool)
	for _, id := range likedPostIDs {
		likedPostMap[id] = true
	}

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
		(SELECT COUNT(*) FROM disliked_posts WHERE post_id = posts.id) AS dislikes,
		(SELECT COUNT(*) FROM comments WHERE post_id = posts.id) AS comments
	FROM
		posts
	INNER JOIN
		users
	ON
		posts.user_id = users.id
	ORDER BY
		posts.created_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		fmt.Printf("Échec de l'exécution de la requête pour obtenir tous les posts. Erreur : %v\n", err)
		return nil, fmt.Errorf("échec de l'exécution de la requête : %v", err)
	}
	defer rows.Close()

	posts := make([]Post, 0)

	for rows.Next() {
		var post Post
		var categoriesStr string

		err := rows.Scan(&post.ID, &post.UserID, &post.Username, &post.Title, &post.Content, &categoriesStr, &post.CreatedAt, &post.Likes, &post.Dislikes, &post.Comments)
		if err != nil {
			fmt.Printf("Échec de la lecture de la ligne lors de la récupération des posts. Erreur : %v\n", err)
			return nil, fmt.Errorf("échec de la lecture de la ligne : %v", err)
		}

		post.Categories = strings.Split(categoriesStr, ",")
		post.IsLiked = likedPostMap[post.ID]
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Échec de la récupération des lignes lors de la récupération des posts. Erreur : %v\n", err)
		return nil, fmt.Errorf("échec de la récupération des lignes : %v", err)
	}

	return posts, nil
}

func GetLikedPostIDs(userID string) ([]string, error) {
	query := "SELECT post_id FROM liked_posts WHERE user_id = ?"
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("échec de l'exécution de la requête : %v", err)
	}
	defer rows.Close()

	var likedPostIDs []string
	for rows.Next() {
		var postID string
		if err := rows.Scan(&postID); err != nil {
			return nil, fmt.Errorf("échec de la lecture de la ligne : %v", err)
		}
		likedPostIDs = append(likedPostIDs, postID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("échec de la récupération des lignes : %v", err)
	}

	return likedPostIDs, nil
}
