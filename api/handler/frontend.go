package handler

import (
	"net/http"
)

func FrontendHandler(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("frontend/dist/"))
	fs.ServeHTTP(w, r)
}
