package handler

import (
	"GoServer/repository"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"slices"
)

type addFriendData struct {
	Action string
	Id     string
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	user, err := authentifacateUser(r)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "access denied", http.StatusUnauthorized)
		return
	}

	var data addFriendData
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	data.Action = html.EscapeString(data.Action)
	data.Id = html.EscapeString(data.Id)

	err = user.AddFriend(data.Id)
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	conversation, err := repository.InsertConversation("DM")
	if err != nil {
		log.Println(err)
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	err = conversation.AddUsers([]string{data.Id, user.Id})
	if err != nil {
		log.Println(err)
		conversation.Delete()
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}

	// fake DM name
	users, err := repository.SelectUsersInConversation(conversation.Id)
	if err != nil {
		log.Println(err)
		conversation.Delete()
		errorResponseJson(w, "error please try agian later", http.StatusBadRequest)
		return
	}
	index := slices.IndexFunc(users, func(u repository.User) bool {
		return u.Id != user.Id
	})
	conversation.Name = users[index].Nickname

	conversationJson, err := json.Marshal(conversation)
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
		"conversation": %s
	}`, conversationJson)))
}
