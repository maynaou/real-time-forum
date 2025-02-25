package models

import (
	"fmt"
	"time"

	database "handler/DataBase/Sqlite"
)

type Comment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
}

func CreateComment(comment Comment, user RegisterRequest) (string, error) {
	query := "INSERT INTO comments (id, user_id, post_id, content, created_at) VALUES (?, ?, ?, ?, ?)"

	_, err := database.DB.Exec(query, comment.ID, user.ID, comment.PostID, comment.Content, comment.CreatedAt)
	if err != nil {
		fmt.Printf("Échec de la création du commentaire : %v\n", err)
		return comment.ID, fmt.Errorf("échec de la création du commentaire : %v", err)
	}

	return comment.ID, nil
}

func GetCommentsByPostID(postID string) ([]Comment, error) {
	query := `
	SELECT comments.id, comments.user_id, comments.post_id, users.nickname, comments.content, comments.created_at,
	(SELECT COUNT(*) FROM liked_posts WHERE post_id = comments.id) AS likes,
	(SELECT COUNT(*) FROM disliked_posts WHERE post_id = comments.id) AS dislikes
	FROM comments
	JOIN users ON comments.user_id = users.id
	WHERE comments.post_id = ?
	ORDER BY comments.created_at DESC
	`

	rows, err := database.DB.Query(query, postID)
	if err != nil {
		fmt.Printf("Échec de l'exécution de la requête : %v\n", err)
		return nil, fmt.Errorf("échec de l'exécution de la requête : %v", err)
	}
	defer rows.Close()

	comments := make([]Comment, 0)
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.PostID, &comment.Username, &comment.Content, &comment.CreatedAt, &comment.Likes, &comment.Dislikes)
		if err != nil {
			fmt.Printf("Échec de la lecture de la ligne : %v\n", err)
			return nil, fmt.Errorf("échec de la lecture de la ligne : %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Échec de la récupération des lignes : %v\n", err)
		return nil, fmt.Errorf("échec de la récupération des lignes : %v", err)
	}

	return comments, nil
}
