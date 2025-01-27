package main

import (
	"fmt"
	"log"
	"net/http"

	database "handler/DataBase"
	userdata "handler/UserData"
	handlers "handler/handlers"
)

var port = "8080"

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.DB.Close()

	http.HandleFunc("/js/", handlers.StaticHandler)
	http.HandleFunc("/styles/", handlers.StaticHandler)
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/api/register", userdata.HandleRegister)
	http.HandleFunc("/api/login", userdata.HandleLogin)
	http.HandleFunc("/api/logout",userdata.HandleLogout)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
