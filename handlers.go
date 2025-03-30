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
	"strconv"
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

func SignUpPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

	tmpl := template.Must(template.ParseFiles("html/auth/Sign_up.html"))
	tmpl.Execute(w, nil)
}

func SignUpGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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

func SignInPostHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		fmt.Fprint(w, "you are allready authenticated")
		return
	}

	tmpl := template.Must(template.ParseFiles("html/auth/Sign_in.html"))
	tmpl.Execute(w, nil)
}

func SignInGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		r.Method = http.MethodGet
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	email := html.EscapeString(r.FormValue("email"))
	password := html.EscapeString(r.FormValue("password"))

	user, err := GetUserLogin(email, password)
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

func MainGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] != nil && session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/home.html"))
	tmpl.Execute(w, nil)
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/messager/home.html"))
	tmpl.Execute(w, nil)
}

func DataConversationsGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := session.Values["user"].(User)

	conversations, err := GetUserConversations(&user)
	if err != nil {
		log.Fatal(err)
	}

	for _, conv := range conversations {
		for _, member := range conv.Members {
			member.Email = ""
			member.Password = ""
		}
	}

	resp, err := json.Marshal(conversations)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp))
	fmt.Fprint(w, string(resp))
}

func DataMessagesGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "server-auth")

	if session.Values["authenticated"] == nil || !session.Values["authenticated"].(bool) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	user := session.Values["user"].(User)
	conversationId, err := strconv.Atoi(html.EscapeString(r.FormValue("id")))
	if err != nil {
		log.Fatal(err)
	}

	if !UserInConversation(user.Id, int64(conversationId)) {
		return
	}

	messages, err := GetConversationMessages(int64(conversationId))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := json.Marshal(messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(string(resp))
	fmt.Fprint(w, string(resp))
}
