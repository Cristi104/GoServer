package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type createConversationData struct {
	Action string
	Name   string
	Users  []string
}

func CreateConversation(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	var data createConversationData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	data.Action = html.EscapeString(data.Action)
	data.Name = html.EscapeString(data.Name)
	for i, v := range data.Users {
		data.Users[i] = html.EscapeString(v)
	}

	conversation, err := repository.InsertConversation(data.Name)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	data.Users = append(data.Users, user.Id)
	err = conversation.AddUsers(data.Users)
	if err != nil {
		log.Println(err)
		conversation.Delete()
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	conversationJson, err := json.Marshal(conversation)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"conversation": %s
	}`, conversationJson)))
}


