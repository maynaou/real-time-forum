package handler

import (
	"encoding/json"
	"fmt"
	models "handler/DataBase/Models"
	utils "handler/Utils"
	"net/http"
)

func Dislike(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createDislike(w, r)
	} else {
		ShowErrorPage(w, "methode not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func createDislike(w http.ResponseWriter, r *http.Request) {
	var dislike models.Liked_Post

	err := json.NewDecoder(r.Body).Decode(&dislike)
	if err != nil {
		fmt.Printf("Failed to decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := utils.GetUserFromSession(r)
	if !ok {
		fmt.Println("User not found in session")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = models.CreateDislike(dislike, user)
	if err != nil {
		fmt.Printf("Failed to create dislike: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(dislike)
	if err != nil {
		fmt.Printf("Failed to marshal dislike data: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		fmt.Printf("Failed to write response data: %v\n", err)
	}
}
