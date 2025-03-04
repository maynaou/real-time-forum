package userdata

import (
	"encoding/json"
	"fmt"
	models "handler/DataBase/Models"
	utils "handler/Utils"
	handler "handler/handlers"
	"net/http"
)

type LoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.ShowErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.ShowErrorPage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	NicknameOREmail := req.Nickname
	if NicknameOREmail == "" {
		NicknameOREmail = req.Email
	}

	user := models.GetUserDetails(req)
	/*if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Println("Invalid credentials")
		return
	}*/

	if (user.Nickname == req.Nickname || user.Email == req.Email) && utils.ComparePasswords(user.Password, req.Password) {
		_, err := utils.SetSession(w, r, user.UserID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Printf("Error setting session: %v", err)
			return
		}

		response := LoginResponse{
			Message:  "Login successful",
			Username: user.Nickname,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Printf("Error encoding JSON: %v", err)
			return
		}

	} else {
		response := LoginResponse{
			Message: "Invalid password or nickname",
		}

		fmt.Println(response)

		w.WriteHeader(http.StatusBadRequest) // Changed to 401 Unauthorized
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			fmt.Printf("Error encoding JSON: %v", err)
			return
		}
	}
}
