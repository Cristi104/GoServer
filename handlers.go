package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := NewUser(username, email, password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, user)
		return
	}

	body, err := os.ReadFile("html/auth/Sign_up.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))

}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := getUserLogin(email, password)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, user)
		return
	}

	body, err := os.ReadFile("html/auth/Sign_in.html")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(body))

}
