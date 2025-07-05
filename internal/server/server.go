package server

import (
	"GoServer/api/handler"
	"net/http"
)

func Run() error {
	http.HandleFunc("/", handler.FrontendHandler)
	return http.ListenAndServe(":8080", nil)
}
