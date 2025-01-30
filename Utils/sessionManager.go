package utils

import (
	"fmt"
	models "handler/DataBase/Models"
	database "handler/DataBase/Sqlite"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func SetSession(w http.ResponseWriter, r *http.Request, id string) (*http.Cookie, error) {
	user := models.RegisterRequest{
		ID: id,
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {

		cookie = &http.Cookie{
			Name:     "session_id",
			Value:    uuid.New().String(),
			Path:     "/",
			HttpOnly: false,
			Expires:  time.Now().Add(2 * time.Hour),
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(w, cookie)

		session := models.Session{
			ID:        cookie.Value,
			UserID:    user.ID,
			ExpiresAt: cookie.Expires,
		}

		_, err = models.CreateSession(session)
		if err != nil {
			fmt.Println("Error creating session:", err)
			return nil, err
		}
	} else {
		session, err := models.GetSessionByID(cookie.Value)
		if err != nil {
			if _, err := DeleteSession(w, r); err != nil {
				fmt.Println("Error deleting session:", err)
				return nil, err
			}

			cookie = &http.Cookie{
				Name:     "session_id",
				Value:    uuid.New().String(),
				Path:     "/",
				HttpOnly: false,
				Expires:  time.Now().Add(2 * time.Hour),
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			}

			http.SetCookie(w, cookie)

			session := models.Session{
				ID:        cookie.Value,
				UserID:    user.ID,
				ExpiresAt: cookie.Expires,
			}

			_, err = models.CreateSession(session)
			if err != nil {
				log.Fatal(err)
			}
		} else if session.UserID != user.ID {
			if _, err := DeleteSession(w, r); err != nil {
				fmt.Println("Error deleting session:", err)
				return nil, err
			}

			cookie = &http.Cookie{
				Name:     "session_id",
				Value:    uuid.New().String(),
				Path:     "/",
				HttpOnly: false,
				Expires:  time.Now().Add(2 * time.Hour),
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			}
			http.SetCookie(w, cookie)

			session := models.Session{
				ID:        cookie.Value,
				UserID:    user.ID,
				ExpiresAt: cookie.Expires,
			}

			_, err = models.CreateSession(session)
			if err != nil {
				fmt.Println("Error creating session:", err)
				log.Fatal(err)
			}
		}
	}

	return cookie, nil
}

func GetUserFromSession(r *http.Request) (models.RegisterRequest, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return models.RegisterRequest{}, false
	}


	database := database.GetDatabaseInstance()
	if database == nil || database.DB == nil {
		fmt.Printf("Database connection error")
		log.Fatal("Database connection error")
		return models.RegisterRequest{}, false
	}

	var userID string
	err = database.DB.QueryRow(`SELECT user_id FROM sessions WHERE id = ? AND expires_at > ?`, cookie.Value, time.Now()).Scan(&userID)
	if err != nil {
		return models.RegisterRequest{}, false
	}

	var user models.RegisterRequest
	err = database.DB.QueryRow(`SELECT id, nickname, first_name, last_name, email, age, password, gender FROM users WHERE id = ?`, userID).Scan(&user.ID, &user.Nickname, &user.FirstName, &user.LastName, &user.Email, &user.Age, &user.Password, &user.Gender)
	if err != nil {
		fmt.Println("User not found in database for the session.")
		return models.RegisterRequest{}, false
	}
	return user, true
}

func DeleteSession(w http.ResponseWriter, r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		fmt.Println("No session cookie found.")
		return nil, nil
	}

	err = models.DeleteSession(cookie.Value)
	if err != nil {
		fmt.Println("Failed to delete session:", err)
		return nil, err
	}

	cookie = &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, cookie)
	return cookie, nil
}
