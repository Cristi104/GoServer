package handler

import (
	"GoServer/repository"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
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

func GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	_, err := authentifacateUser(r)
	if err != nil {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()
	username := query.Get("username")

	users, err := repository.SelectUsersByUsername(username)
	if err != nil {
		errorResponseJson(w, "the user doesn't exist", http.StatusBadRequest)
		return
	}

	jsonString := strings.Builder{}
	for i, u := range users {
		jsonString.Write([]byte(fmt.Sprintf(`{"id": "%s", "username": "%s", "nickname": "%s" }`, u.Id, u.Username, u.Nickname)))
		if i != len(users)-1 {
			jsonString.Write([]byte(","))
		}
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"success": true, "error": "none", "profiles": [%s] }`, jsonString.String())))
}
