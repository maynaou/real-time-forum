package handler

import "net/http"

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ShowErrorPage(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		ShowErrorPage(w, "Page not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "templates/HomePage.html")
}
