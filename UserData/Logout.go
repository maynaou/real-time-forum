package userdata

import (
	"encoding/json"
	database "handler/DataBase/Sqlite"
	"net/http"
	"time"
)

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No active session", http.StatusBadRequest)
		return
	}

	_, err = database.GetDatabaseInstance().DB.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-1 * time.Hour),
	})

	response := JsonResponse{Message: "Logout successful"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
