package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

		err = ValidateSession(session.Value)
		if err != nil {
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

		err = ValidateSession(session.Value)
		if err != nil {
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
