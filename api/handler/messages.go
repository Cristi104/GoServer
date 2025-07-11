package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type createMessageData struct {
	ConversationId string
	Body           string
}

type listener struct {
	flusher http.Flusher
	w       http.ResponseWriter
}

var messageListeners = make(map[string][]listener)

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

	for _, f := range messageListeners[data.ConversationId] {
		f.w.Write([]byte(fmt.Sprintf(`{
			success": true, 
			error": "none", 
			message": [%s]
		}`, messageJson)))
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"message": %s
	}`, messageJson)))
}

func GetAllMessages(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	conversationId := chi.URLParam(r, "id")
	messages, err := repository.SelectMessagesByConversation(conversationId, user.Id)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	messagesJson, err := json.Marshal(messages)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"message": %s
	}`, messagesJson)))
}

func MessageListener(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	conversationId := chi.URLParam(r, "id")
	messages, err := repository.SelectMessagesByConversation(conversationId, user.Id)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	messagesJson, err := json.Marshal(messages)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Write([]byte(fmt.Sprintf(`{
		"success": true, 
		"error": "none", 
		"message": %s
	}`, messagesJson)))

	flusher, ok := w.(http.Flusher)
	if !ok {
		errorResponseJson(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	messageListeners[conversationId] = append(messageListeners[conversationId], listener{w: w, flusher: flusher})
}
