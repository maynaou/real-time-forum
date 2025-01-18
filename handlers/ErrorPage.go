package handler

import (
	"net/http"
	"text/template"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func ShowErrorPage(w http.ResponseWriter, errMsg string, statusCode int) {
	data := struct {
		Err string
	}{
		Err: errMsg,
	}
	w.WriteHeader(statusCode)
	if err := templates.ExecuteTemplate(w, "ErrorPage.html", data); err != nil {
		http.Error(w, "Error displaying the error page", http.StatusInternalServerError)
	}
}
