package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	auth, err := r.Cookie("auth")
	if err != nil {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	token, err := validateJWT(auth.Value)
	if err != nil {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	var user repository.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp := claims["exp"].(float64)
		expTime := time.Unix(int64(exp), 0)
		if expTime.Before(time.Now()) {
			errorResponseJson(w, "access denied", http.StatusUnauthorized)
			return
		}

		userData := claims["userData"].(string)
		err = json.Unmarshal([]byte(userData), &user)
		if err != nil {
			errorResponseJson(w, "access denied", http.StatusUnauthorized)
			return
		}
	} else {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	id := query.Get("id")

	if id == "" {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"success": true, 
			"error": "none", 
			"profile": {
				"id": "%s", 
				"username": "%s", 
				"nickname": "%s", 
				"email": "%s"
			}
		}`, user.Id, user.Username, user.Nickname, user.Email)))
	} else {
		user, err = repository.SelectUserById(id)
		if err != nil {
			errorResponseJson(w, "the user doesn't exist", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{
			"success": true, 
			"error": "none", 
			"profile": {
				"id": "%s", 
				"username": "%s", 
				"nickname": "%s", 
			}
		}`, user.Id, user.Username, user.Nickname)))
	}
}
