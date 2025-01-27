package userdata

import (
	"encoding/json"
	"fmt"
	"net/http"

	database "handler/DataBase"
	handler "handler/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type JsonResponse struct {
	Message string `json:"message"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("jjjjj2")
	if r.Method != http.MethodPost {
		handler.ShowErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.ShowErrorPage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(req)

	if req.Nickname == "" || req.Email == "" || req.Password == "" {
		handler.ShowErrorPage(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ? OR email = ?)",
		req.Nickname, req.Email).Scan(&exists)
	if err != nil {
		handler.ShowErrorPage(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		handler.ShowErrorPage(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		handler.ShowErrorPage(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	userID := uuid.New().String()

	_, err = database.DB.Exec(`
        INSERT INTO users (id,nickname, email, password, first_name, last_name, age, gender)
        VALUES (?,?, ?, ?, ?, ?, ?, ?)`,
		userID, req.Nickname, req.Email, hashedPass, req.FirstName, req.LastName, req.Age, req.Gender)
	if err != nil {
		handler.ShowErrorPage(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	response := JsonResponse{Message: "User created successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
