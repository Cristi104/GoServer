package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type signInCredentials struct {
	Email    string
	Password string
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var data signInCredentials

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data.Email = html.EscapeString(data.Email)
	data.Password = html.EscapeString(data.Password)

	user, err := repository.SelectUserBySignIn(data.Email, data.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusBadRequest)
		return
	}

	userData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userData":  userData,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	signedToken, err := token.SignedString([]byte(JWTKey))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-type", "application/json")

	cookie := http.Cookie{Name: "auth", Value: signedToken, Quoted: false, MaxAge: 60 * 60 * 1, Secure: false, HttpOnly: false, SameSite: http.SameSiteStrictMode}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"success\":%s, \"error\":\"%s\"}", "true", "none")))
}
