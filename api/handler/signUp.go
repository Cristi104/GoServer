package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"regexp"
	"unicode"
)

func isValidPassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func isValidUsername(username string) bool {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
	)

	if len(username) >= 8 {
		hasMinLen = true
	}

	for _, ch := range username {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		}
	}
	return hasMinLen && (hasUpper || hasLower)
}

type signUpCredentials struct {
	Username string
	Email    string
	Password string
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var data signUpCredentials

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	data.Username = html.EscapeString(data.Username)
	data.Email = html.EscapeString(data.Email)
	data.Password = html.EscapeString(data.Password)

	if !isValidUsername(data.Username) {
		errorResponseJson(w, "invalid username (at least 8 letters long)", http.StatusBadRequest)
		return
	}

	matched, err := regexp.Match("^[A-Za-z0-9]+@[a-z]+\\.[a-z][a-z][a-z]?$", []byte(data.Email))
	if err != nil || !matched {
		errorResponseJson(w, "invalid email", http.StatusBadRequest)
		return
	}

	if !isValidPassword(data.Password) {
		errorResponseJson(w, "invalid password (atleast 8 characters long, one uppercase, one lowercase, one letter and one symbol)", http.StatusBadRequest)
		return
	}

	user, err := repository.InsertUser(data.Username, data.Email, data.Password)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "email or username allready in use", http.StatusBadRequest)
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
