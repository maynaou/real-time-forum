package handler

import (
	"encoding/json"
	"net/http"

	models "handler/DataBase/Models"
)

func Message(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getMessage(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	sender := r.URL.Query().Get("sender")
	receiver := r.URL.Query().Get("receiver")
	before := r.URL.Query().Get("before")
	if sender == "" || receiver == "" {
		http.Error(w, "Sender and Receiver are required", http.StatusBadRequest)
		return
	}

	messages, err := models.GetMessages(sender, receiver, before)
	if err != nil {
		http.Error(w, "Failed to fetch messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
