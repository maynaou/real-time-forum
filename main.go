package main

import (
	"fmt"
	"net/http"

	handlers "handler/handlers"
)

var port = "8080"

func main() {
	http.HandleFunc("/js/", handlers.StaticHandler)
	http.HandleFunc("/styles/", handlers.StaticHandler)
	http.HandleFunc("/", handlers.HomePage)
	fmt.Println("Server started on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
