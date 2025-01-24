package userdata

import (
	"encoding/json"
	"net/http"
	"time"

	database "handler/DataBase"
	handler "handler/handlers"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.ShowErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.ShowErrorPage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user struct {
		ID       string
		Password string
	}

	err := database.DB.QueryRow(`
	SELECT id, password FROM users WHERE email = ? OR nickname = ?
	`, req.Login, req.Login).Scan(&user.ID, &user.Password) // Correction ici
	if err != nil {
		handler.ShowErrorPage(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = database.DB.Exec(`
        INSERT INTO sessions (id, user_id, expires_at) 
        VALUES (?, ?, ?)`,
		sessionID, user.ID, expiresAt)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: expiresAt,
	})
}
