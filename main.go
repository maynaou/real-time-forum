package main

import (
	"fmt"
	"log"
	"net/http"

	database "handler/DataBase/Sqlite"
	userdata "handler/UserData"
	utils "handler/Utils"
	handler "handler/handlers"
)

var port = "8075"

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
	http.HandleFunc("/api/post", utils.Middleware(handler.Post))
	http.HandleFunc("/api/like", utils.Middleware(handler.Like))
	http.HandleFunc("/api/dislike", utils.Middleware(handler.Dislike))
	http.HandleFunc("/api/comment", utils.Middleware(handler.Comment))
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
