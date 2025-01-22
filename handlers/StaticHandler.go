package handler

import (
	"net/http"
	"os"
	"strings"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/js/") {
		fs := http.StripPrefix("/js/", http.FileServer(http.Dir("js")))
		_, err := os.Stat("." + r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/") || err != nil {
			ShowErrorPage(w, "don't have access", http.StatusForbidden)
			return
		}
		fs.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/styles/") {
		fs := http.StripPrefix("/styles/", http.FileServer(http.Dir("styles/")))
		_, err := os.Stat("." + r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/") || err != nil {
			ShowErrorPage(w, "don't have access", http.StatusForbidden)
			return
		}
		fs.ServeHTTP(w, r)
	}
}
