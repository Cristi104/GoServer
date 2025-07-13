package handler

import (
	"GoServer/repository"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetAllConversationUsers(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	users, err := repository.SelectUsersInConversation(id)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	if !slices.ContainsFunc(users, func(u repository.User) bool {
		return u.Id == user.Id
	}) {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	var jsonString strings.Builder
	for i, u := range users {
		fmt.Fprintf(&jsonString, `{"id": "%s", "nickname": "%s"}`, u.Id, u.Nickname)
		if i != len(users)-1 {
			fmt.Fprint(&jsonString, ",")
		}
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"users": [%s]
	}`, jsonString.String())))
}
