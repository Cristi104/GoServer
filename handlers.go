package main

import (
	"encoding/gob"
	"encoding/json"
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

func signUpPageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

	tmpl := template.Must(template.ParseFiles("html/auth/Sign_up.html"))
	tmpl.Execute(w, nil)
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

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
}

func signInPageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

	tmpl := template.Must(template.ParseFiles("html/auth/Sign_in.html"))
	tmpl.Execute(w, nil)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

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

	r.Method = http.MethodGet
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/home.html"))
	tmpl.Execute(w, nil)
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/messager/main.html"))
	tmpl.Execute(w, nil)
}

func homePageLoader(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := session.Values["user"].(User)

	friends, err := getFriends(&user)
	if err != nil {
		log.Fatal(err)
	}

	var usernames []string
	for _, friend := range friends {
		usernames = append(usernames, friend.Username)
	}

	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(usernames)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, string(resp))
}
