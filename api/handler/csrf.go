package handler

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.WriteHeader(http.StatusOK)
}
