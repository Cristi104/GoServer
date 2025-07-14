package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
)

var messageListeners = make(map[string][]chan string)

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

	if string(messagesJson) == "null" {
		messagesJson = []byte("[]")
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(w, `data: {"success": true, "error": "none", "messages": %s}%s`, messagesJson, "\n\n")

	flusher, ok := w.(http.Flusher)
	if !ok {
		errorResponseJson(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	notify := r.Context().Done()
	messageChan := make(chan string)
	messageListeners[conversationId] = append(messageListeners[conversationId], messageChan)
	for {
		select {
		case <-notify:
			messageListeners[conversationId] = slices.DeleteFunc(messageListeners[conversationId], func(ch chan string) bool {
				return ch == messageChan
			})
			return
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}

}
