package main

import (
	"fmt"
	"log"
	"net/http"

	database "handler/DataBase/Sqlite"
	userdata "handler/UserData"
	handler "handler/handlers"
)

var port = "8080"

func main() {

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	db := database.GetDatabaseInstance()
	if db.DB != nil {
		log.Println("Connexion à la base de données réussie.")
	}

	http.HandleFunc("/js/", handler.StaticHandler)
	http.HandleFunc("/styles/", handler.StaticHandler)
	http.HandleFunc("/", handler.HomePage)
	http.HandleFunc("/api/register", userdata.HandleRegister)
	http.HandleFunc("/api/login", userdata.HandleLogin)
	http.HandleFunc("/api/logout", userdata.HandleLogout)
	http.HandleFunc("/api/post", handler.Post)
	http.HandleFunc("/api/like", handler.Like)
	http.HandleFunc("/api/dislike", handler.Dislike)
	http.HandleFunc("/api/comments", handler.Comment)
	http.HandleFunc("/api/comment",handler.Comment)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
