package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetAllConversations(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	conversations, err := repository.SelectConversationsByUser(user.Id)
	conversationsJson, err := json.Marshal(conversations)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	if string(conversationsJson) == "null" {
		conversationsJson = []byte("[]")
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"conversations": %s
	}`, conversationsJson)))
}

