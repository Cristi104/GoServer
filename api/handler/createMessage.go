package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type createMessageData struct {
	ConversationId string
	Body           string
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	var data createMessageData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	data.Body = html.EscapeString(data.Body)
	data.ConversationId = html.EscapeString(data.ConversationId)

	message, err := repository.InsertMessage(data.Body, user.Id, data.ConversationId)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	for _, ch := range messageListeners[data.ConversationId] {
		ch <- fmt.Sprintf(`{"success": true, "error": "none", "messages": [%s]}`, messageJson)
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"message": %s
	}`, messageJson)))
}



