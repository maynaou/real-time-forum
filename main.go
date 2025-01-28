package main

import (
	"fmt"
	"log"
	"net/http"

	database "handler/DataBase/Sqlite"
	userdata "handler/UserData"
	handler "handler/handlers"
)

var port = "8088"

func main() {
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
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
