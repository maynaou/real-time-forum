package utils

import (
	"encoding/json"
	"fmt"
	database "handler/DataBase/Sqlite"
	"log"
	"net/http"
	"time"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err != nil || session.Value == "" {
			log.Println("Unauthorized. Redirecting to login.")
			w.WriteHeader(http.StatusUnauthorized) // Return 401 status
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		next(w, r)
	}
}

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session_id")
		if err != nil || session.Value == "" {
			next(w, r)
			return
		}

		fmt.Println(session)

		database := database.GetDatabaseInstance()
		if database == nil || database.DB == nil {
			fmt.Printf("Database connection error")
			log.Fatal("Database connection error")
			return
		}

		var expiresAt time.Time
		err = database.DB.QueryRow(
			"SELECT expires_at FROM sessions WHERE id = ?;",
			session.Value,
		).Scan(&expiresAt)

		if err != nil || time.Now().After(expiresAt) {
			// Si la session est invalide ou expirée, passer au prochain handler
			next(w, r)
			return
		}

		data, err := json.Marshal(map[string]string{"authenticated": "true"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(data); err != nil {
			fmt.Printf("failed to write response data: %v", err)
		}
	}
}
