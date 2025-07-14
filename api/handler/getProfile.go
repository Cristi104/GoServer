package handler

import (
	"GoServer/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	if id == user.Id {
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
			log.Println(err)
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
