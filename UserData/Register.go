package userdata

import (
	"encoding/json"
	"fmt"
	models "handler/DataBase/Models"
	utils "handler/Utils"
	handler "handler/handlers"
	"net/http"

	"github.com/google/uuid"
)

type JsonResponse struct {
	Message string `json:"message"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		handler.ShowErrorPage(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	req.ID = uuid.New().String()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handler.ShowErrorPage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validationErrors := utils.ValidateRegisterFornData(req)
	fmt.Println(validationErrors)
	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(validationErrors); err != nil {
			fmt.Printf("Failed to encode validation errors: %v", err)
		}
		return
	}

	hashedPass, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPass)
	_, err = models.CreateUser(req)

	if err != nil {
		handler.ShowErrorPage(w, "Database error", http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errorMsg := map[string]string{"error": "Username or email already exists."}
		if err := json.NewEncoder(w).Encode(errorMsg); err != nil {
			fmt.Printf("Failed to encode error message: %v", err)
		}
		return
	}

	response := JsonResponse{
		Message: "register successful",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
		fmt.Printf("Error encoding JSON: %v", err)
		return
	}
}
