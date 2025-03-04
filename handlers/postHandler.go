package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "handler/DataBase/Models"
	utils "handler/Utils"

	"github.com/google/uuid"
)

func Post(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createPost(w, r)
	} else if r.Method == http.MethodGet {
		getPosts(w, r)
	} else {
		ShowErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	post.ID = uuid.New().String()

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		fmt.Printf("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, ok := utils.GetUserFromSession(r)

	if !ok {
		fmt.Println("user not found in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validationErrors := utils.ValidationPostInput(post)

	if len(validationErrors) > 0 {
		fmt.Printf("Validation errors: %v", validationErrors)
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validationErrors); err != nil {
			fmt.Printf("faild to encode validation errors: %v", err)
		}
		return
	}

	post.ID, err = models.CreatePost(post, user)
	if err != nil {
		fmt.Printf("failed to marshal post data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Retourner le post créé avec son ID
	response := map[string]interface{}{
		"id":         post.ID,
		"username":   user.Nickname, // Assurez-vous de renvoyer le bon username
		"title":      post.Title,
		"content":    post.Content,
		"category":   post.Categories,
		"created_at": post.CreatedAt,
		"comments":   0,
		"likes":      0,
		"dislikes":   0,
		"message":    "post successful",
	}

	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data) // Envoyer le post créé en réponse
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	posts, err := models.GetAllPosts(user.ID)
	if err != nil {
		fmt.Printf("Failed to marshal posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(posts)
	if err != nil {
		fmt.Printf("Failed to marshal posts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		fmt.Printf("failed to write response data: %v", err)
	}
}
