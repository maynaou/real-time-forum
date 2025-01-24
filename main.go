package main

import (
	"fmt"
	"net/http"

	userdata "handler/UserData"
	handlers "handler/handlers"
)

var port = "8080"

func main() {
	http.HandleFunc("/js/", handlers.StaticHandler)
	http.HandleFunc("/styles/", handlers.StaticHandler)
	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/api/register", userdata.HandleRegister)
	http.HandleFunc("/api/login", userdata.HandleLogin)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
