package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

func ConversationsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	conversations, err :=repository.SelectConversationsByUser(user.Id)
	conversationsJson, err := json.Marshal(conversations)
	if err != nil {
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"conversations": %s
	}`, conversationsJson)))
}
