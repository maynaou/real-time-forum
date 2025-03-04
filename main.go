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

var port = "8066"

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/js/", handler.StaticHandler)
	http.HandleFunc("/styles/", handler.StaticHandler)
	http.HandleFunc("/", handler.HomePage)
	http.HandleFunc("/api/register", utils.Auth(userdata.HandleRegister))
	http.HandleFunc("/api/login", utils.Auth(userdata.HandleLogin))
	http.HandleFunc("/api/logout", userdata.HandleLogout)
	http.HandleFunc("/api/user", handler.User)
	http.HandleFunc("/api/post/", utils.Middleware(handler.Post))
	http.HandleFunc("/api/like", utils.Middleware(handler.Like))
	http.HandleFunc("/api/dislike", utils.Middleware(handler.Dislike))
	http.HandleFunc("/api/comment/", utils.Middleware(handler.Comment))
	http.HandleFunc("/api/message", utils.Middleware(handler.Message))
	http.HandleFunc("/ws", handler.WebSocket)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
