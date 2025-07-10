package repository

import (
	"fmt"
	"strings"
)

const INSERT_CONVERSATION_SQL = "INSERT INTO conversations(name) VALUES($1) RETURNING id, create_date;"
const SELECT_CONVERSATION_BY_ID_SQL = "SELECT * FROM conversations WHERE id = $1;"
// const SELECT_CONVERSATIONS_BY_USER_SQL = "SELECT c.* FROM conversations c LEFT JOIN in_conversation ic ON c.id = ic.conversation_id WHERE ic.user_id = $1;"
const SELECT_CONVERSATIONS_BY_USER_SQL = `
SELECT c.id,
CASE
	WHEN c.name = 'DM' THEN (
		SELECT u.nickname 
		FROM users u LEFT JOIN in_conversation ic2 on u.id = ic2.user_id 
		WHERE ic2.conversation_id = c.id AND ic2.user_id != $1
	)
	ELSE c.name
END AS "name", c.create_date
FROM conversations c LEFT JOIN in_conversation ic on c.id = ic.conversation_id
WHERE ic.user_id = '7be3e465-e058-49a2-9126-48af79b63cea';
`
const UPDATE_CONVERSATION_SQL = "UPDATE conversations SET name = $1 WHERE id = $2;"
const DELETE_CONVERSATION_SQL = "DELETE FROM conversations WHERE id = $1;"
const INSERT_USERS_IN_CONVERSATION_SQL = "INSERT INTO in_conversation(conversation_id, user_id) VALUES"

type Conversation struct {
	Id         string
	Name       string
	CreateDate string
}

func InsertConversation(name string) (Conversation, error) {
	var conversation Conversation

	conversation.Name = name
	err := DatabaseConnection.QueryRow(INSERT_CONVERSATION_SQL, name).Scan(&conversation.Id, &conversation.CreateDate)
	if err != nil {
		return Conversation{}, err
	}

	return conversation, nil
}

func SelectConversationById(id string) (Conversation, error) {
	var conversation Conversation

	err := DatabaseConnection.QueryRow(SELECT_CONVERSATION_BY_ID_SQL, id).Scan(&conversation.Id, &conversation.Name, &conversation.CreateDate)
	if err != nil {
		return Conversation{}, err
	}

	return conversation, nil
}

func SelectConversationsByUser(userId string) ([]Conversation, error) {
	rows, err := DatabaseConnection.Query(SELECT_CONVERSATIONS_BY_USER_SQL, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	var conversation Conversation
	for rows.Next() {
		err := rows.Scan(&conversation.Id, &conversation.Name, &conversation.CreateDate)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func (c *Conversation) AddUsers(userIds []string) error {
	var buffer strings.Builder
	buffer.WriteString(INSERT_USERS_IN_CONVERSATION_SQL)

	for i := range userIds {
		buffer.WriteString(fmt.Sprintf(" ($1, $%d)", i+2))

		if i != len(userIds)-1 {
			buffer.WriteString(",")
		} else {
			buffer.WriteString(";")
		}
	}

	args := make([]any, len(userIds)+1)
	args[0] = c.Id
	for i, v := range userIds {
		args[i+1] = v
	}

	sqlString := buffer.String()
	_, err := DatabaseConnection.Exec(sqlString, args...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) AddUser(userId string) error {
	var buffer strings.Builder

	buffer.WriteString(INSERT_USERS_IN_CONVERSATION_SQL)
	buffer.WriteString(" ($1, $2);")
	sqlString := buffer.String()

	_, err := DatabaseConnection.Exec(sqlString, c.Id, userId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) Update() error {
	_, err := DatabaseConnection.Exec(UPDATE_CONVERSATION_SQL, c.Name, c.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) Delete() error {
	_, err := DatabaseConnection.Exec(DELETE_CONVERSATION_SQL, c.Id)
	if err != nil {
		return err
	}
	return nil
}

