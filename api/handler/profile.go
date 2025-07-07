package handler

import (
	"GoServer/repository"
	"fmt"
	"net/http"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
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
