package handler

import (
	"net/http"
)

func Like(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createLike(w, r)
	} else {
		ShowErrorPage(w, "methode not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func Likes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getLikes(w, r)
	} else {
		ShowErrorPage(w, "methode not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func createLike(w http.ResponseWriter, r *http.Request) {

}

func getLikes(w http.ResponseWriter, r *http.Request) {

}
