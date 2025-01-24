package userdata

import (
	"encoding/json"
	"net/http"

	database "handler/DataBase"
	handler "handler/handlers"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       string `json:"age"`
	Gender    string `json:"gender"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.ShowErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.ShowErrorPage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validation des champs requis
	if req.Nickname == "" || req.Email == "" || req.Password == "" {
		handler.ShowErrorPage(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Vérification des doublons d'utilisateur
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

	// Hachage du mot de passe
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		handler.ShowErrorPage(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insertion de l'utilisateur dans la base de données
	_, err = database.DB.Exec(`
        INSERT INTO users (nickname, email, password, first_name, last_name, age, gender)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.Nickname, req.Email, hashedPass, req.FirstName, req.LastName, req.Age, req.Gender)
	if err != nil {
		handler.ShowErrorPage(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Réponse de succès
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
