package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	models "handler/DataBase/Models"
	utils "handler/Utils"
	"net/http"

	"github.com/google/uuid"
)

func Comment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createComment(w, r)
	} else if r.Method == http.MethodGet {
		getComment(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	comment.ID = uuid.New().String()

	user, ok := utils.GetUserFromSession(r)
	if !ok {
		fmt.Println("User not found in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		fmt.Printf("Failed to decode comment from request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validationErrors := utils.ValidateCommentInput(comment)
	if len(validationErrors) > 0 {
		fmt.Println("Comment input validation failed")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validationErrors); err != nil {
			fmt.Printf("Failed to encode validation errors: %v", err)
		}
		return
	}

	_, err = models.CreateComment(comment, user)
	if err != nil {
		fmt.Printf("Failed to create comment: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(comment)
	if err != nil {
		fmt.Printf("Failed to marshal comment: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		fmt.Printf("Failed to write response data: %v", err)
	}
}

func getComment(w http.ResponseWriter, r *http.Request) {
	postID := r.Header.Get("X-Requested-With")

	comment, err := models.GetCommentsByPostID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No comments found for post ID %s", postID)
			http.NotFound(w, r)
			return
		}
		fmt.Printf("Failed to fetch comments for post ID %s: %v", postID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(comment)
	if err != nil {
		fmt.Printf("Failed to marshal comments: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		fmt.Printf("Failed to write response data: %v", err)
	}
}
