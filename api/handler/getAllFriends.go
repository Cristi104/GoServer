package handler

import (
	"GoServer/repository"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func GetAllFriends(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	friends, err := repository.SelectUserFriends(user.Id)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try again later", http.StatusUnauthorized)
		return
	}

	jsonString := strings.Builder{}
	for i, u := range friends {
		jsonString.Write([]byte(fmt.Sprintf(`{"id": "%s", "username": "%s", "nickname": "%s" }`, u.Id, u.Username, u.Nickname)))
		if i != len(friends)-1 {
			jsonString.Write([]byte(","))
		}
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"friends": [%s]
	}`, jsonString.String())))
}
