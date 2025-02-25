package models

import (
	"log"

	database "handler/DataBase/Sqlite"
)

type Liked_Post struct {
	Post_ID string `json:"post_id"`
	User_ID string `json:"user_id"`
}

func CreateLike(like Liked_Post, user RegisterRequest) error {
	// Vérifier si le like existe déjà
	var likeCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID).Scan(&likeCount)
	if err != nil {
		log.Printf("Échec de la vérification de l'existence du like : %v\n", err)
		return err
	}

	// Si le like n'existe pas, vérifier pour un dislike
	if likeCount == 0 {
		// Vérifier si un dislike existe
		var dislikeCount int
		err = database.DB.QueryRow("SELECT COUNT(*) FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID).Scan(&dislikeCount)
		if err != nil {
			log.Printf("Échec de la vérification de l'existence du dislike : %v\n", err)
			return err
		}

		// Si un dislike existe, le supprimer
		if dislikeCount > 0 {
			_, err = database.DB.Exec("DELETE FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID)
			if err != nil {
				log.Printf("Échec de la suppression du dislike : %v\n", err)
				return err
			}
		}

		// Insérer le like
		_, err = database.DB.Exec("INSERT INTO liked_posts (user_id, post_id) VALUES (?, ?)", user.ID, like.Post_ID)
		if err != nil {
			log.Printf("Échec de l'insertion du like : %v\n", err)
			return err
		}
	} else {
		// Si le like existe, le supprimer
		_, err = database.DB.Exec("DELETE FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, like.Post_ID)
		if err != nil {
			log.Printf("Échec de la suppression du like : %v\n", err)
			return err
		}
	}

	return nil
}

func CountLikes(postID string) (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM liked_posts WHERE post_id = ?", postID).Scan(&count)
	return count, err
}
