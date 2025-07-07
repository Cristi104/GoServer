package server

import (
	"GoServer/api/handler"
	"net/http"
)

func Run() error {
	http.HandleFunc("/", handler.FrontendHandler)
	http.HandleFunc("/api/auth/signin", handler.SignInHandler)
	http.HandleFunc("/api/auth/signup", handler.SignUpHandler)
	return http.ListenAndServe(":8080", nil)
}
