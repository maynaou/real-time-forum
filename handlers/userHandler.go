package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "handler/DataBase/Models"
	utils "handler/Utils"
)

func User(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getUser(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Users(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getUsers(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromSession(r)

	if !ok {
		fmt.Println("user not found in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Failed to marshal user data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		fmt.Printf("Failed to write response data: %v", err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	user, ok := utils.GetUserFromSession(r)

	if !ok {
		fmt.Println("user not found in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	onlineUsers := GetActiveUsers(w, r)

	onlineMap := make(map[string]bool)
	for _, id := range onlineUsers {
		onlineMap[id] = true
	}

	users, err := models.GetAllUsers(user.ID, onlineMap)
	if err != nil {
		fmt.Printf("Failed to retrieve all users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(users)
	if err != nil {
		fmt.Printf("Failed to marshal users data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(data); err != nil {
		fmt.Printf("Failed to write response data: %v", err)
	}
}
