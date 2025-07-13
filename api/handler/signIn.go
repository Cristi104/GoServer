package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type signInCredentials struct {
	Email    string
	Password string
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var data signInCredentials

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	data.Email = html.EscapeString(data.Email)
	data.Password = html.EscapeString(data.Password)

	user, err := repository.SelectUserBySignIn(data.Email, data.Password)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	userData, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	signedToken, err := createAuthJWT(string(userData))
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-type", "application/json")

	cookie := http.Cookie{Name: "auth", Value: signedToken, Path: "/", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success":true, "error":"none", "userId": "%s"}`, user.Id)))
}
