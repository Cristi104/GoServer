package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/signup/", signUpHandler)
	http.HandleFunc("/signin/", signInHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
