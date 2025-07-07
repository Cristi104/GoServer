package server

import (
	"GoServer/api/handler"
	"net/http"
)

func Run() error {
	http.HandleFunc("/", handler.FrontendHandler)
	http.HandleFunc("POST /api/auth/signin", handler.SignInHandler)
	http.HandleFunc("POST /api/auth/signup", handler.SignUpHandler)
	http.HandleFunc("GET /api/profile", handler.ProfileHandler)
	http.HandleFunc("GET /api/conversations", handler.ConversationsHandler)
	return http.ListenAndServe(":8080", nil)
}
