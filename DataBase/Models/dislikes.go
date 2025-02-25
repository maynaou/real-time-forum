package models

import (
	"log"

	database "handler/DataBase/Sqlite"
)

func CreateDislike(dislike Liked_Post, user RegisterRequest) error {
	// Vérifier si le dislike existe déjà
	var dislikeCount int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID).Scan(&dislikeCount)
	if err != nil {
		log.Printf("Échec de la vérification de l'existence du dislike : %v\n", err)
		return err
	}

	// Si le dislike n'existe pas, vérifier pour un like
	if dislikeCount == 0 {
		// Vérifier si un like existe
		var likeCount int
		err = database.DB.QueryRow("SELECT COUNT(*) FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID).Scan(&likeCount)
		if err != nil {
			log.Printf("Échec de la vérification de l'existence du like : %v\n", err)
			return err
		}

		// Si un like existe, le supprimer
		if likeCount > 0 {
			_, err = database.DB.Exec("DELETE FROM liked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID)
			if err != nil {
				log.Printf("Échec de la suppression du like : %v\n", err)
				return err
			}
		}

		// Insérer le dislike
		_, err = database.DB.Exec("INSERT INTO disliked_posts (user_id, post_id) VALUES (?, ?)", user.ID, dislike.Post_ID)
		if err != nil {
			log.Printf("Échec de l'insertion du dislike : %v\n", err)
			return err
		}
	} else {
		// Si le dislike existe, le supprimer
		_, err = database.DB.Exec("DELETE FROM disliked_posts WHERE user_id = ? AND post_id = ?", user.ID, dislike.Post_ID)
		if err != nil {
			log.Printf("Échec de la suppression du dislike : %v\n", err)
			return err
		}
	}

	return nil
}

func CountDislikes(postID string) (int, error) {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM disliked_posts WHERE post_id = ?", postID).Scan(&count)
	return count, err
}
