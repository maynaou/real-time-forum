package main

import (
	"fmt"
	"log"
	"net/http"

	database "handler/DataBase/Sqlite"
	userdata "handler/UserData"
	handler "handler/handlers"
)

var port = "8083"

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
	http.HandleFunc("/api/posts", handler.Posts)
	http.HandleFunc("/api/post", handler.Post)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
