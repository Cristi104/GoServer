package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	gob.Register(User{})

	key, err := os.ReadFile("config/session_key.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}

		log.Print(err)
		log.Print("new random key is generated")

		key = securecookie.GenerateRandomKey(32)

		err = os.WriteFile("config/session_key.txt", key, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
	store = sessions.NewCookieStore(key)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == true {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

	if r.Method == http.MethodPost {
		username := html.EscapeString(r.FormValue("username"))
		email := html.EscapeString(r.FormValue("email"))
		password := html.EscapeString(r.FormValue("password"))

		user, err := NewUser(username, email, password)
		if err != nil {
			log.Fatal(err)
		}

		session.Values["authenticated"] = true
		session.Values["user"] = *user
		err = session.Save(r, w)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, user)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/auth/Sign_up.html"))
	tmpl.Execute(w, nil)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == true {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

	if r.Method == http.MethodPost {
		email := html.EscapeString(r.FormValue("email"))
		password := html.EscapeString(r.FormValue("password"))

		user, err := getUserLogin(email, password)
		if err != nil {
			log.Fatal(err)
		}

		session.Values["authenticated"] = true
		session.Values["user"] = *user
		err = session.Save(r, w)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Fprint(w, user)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/auth/Sign_in.html"))
	tmpl.Execute(w, nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == true {
		http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/home.html"))
	tmpl.Execute(w, nil)
}
