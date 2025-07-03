package repository

import (
	"fmt"
	"strings"
)

const INSERT_CONVERSATION_SQL = "INSERT INTO conversations(name) VALUES($1) RETURNING id, create_date;"
const SELECT_CONVERSATION_BY_ID_SQL = "SELECT * FROM conversations WHERE id = $1;"
const SELECT_CONVERSATIONS_BY_USER_SQL = "SELECT c.* FROM conversations c LEFT JOIN in_conversation ic ON c.id = ic.conversation_id WHERE ic.user_id = $1;"
const UPDATE_CONVERSATION_SQL = "UPDATE conversations SET name = $1 WHERE id = $2;"
const DELETE_CONVERSATION_SQL = "DELETE FROM conversations WHERE id = $1;"
const INSERT_USERS_IN_CONVERSATION_SQL = "INSERT INTO in_conversation(conversation_id, user_id) VALUES"

type Conversation struct {
	id         string
	name       string
	createDate string
}

func InsertConversation(name string) (*Conversation, error) {
	var conversation Conversation

	conversation.name = name
	err := DatabaseConnection.QueryRow(INSERT_CONVERSATION_SQL, name).Scan(&conversation.id, &conversation.createDate)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func SelectConversationById(id string) (*Conversation, error) {
	var conversation Conversation

	err := DatabaseConnection.QueryRow(SELECT_CONVERSATION_BY_ID_SQL, id).Scan(&conversation.id, &conversation.name, &conversation.createDate)
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func SelectConversationsByUser(userId string) ([]*Conversation, error) {
	var conversations []*Conversation

	rows, err := DatabaseConnection.Query(SELECT_CONVERSATIONS_BY_USER_SQL, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversation Conversation
	for rows.Next() {
		err := rows.Scan(&conversation.id, &conversation.name, &conversation.createDate)
		if err != nil {
			return nil, err
		}

		copy := conversation
		conversations = append(conversations, &copy)
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
	args[0] = c.Id()
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

	_, err := DatabaseConnection.Exec(sqlString, c.id, userId)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) Update() error {
	_, err := DatabaseConnection.Exec(UPDATE_CONVERSATION_SQL, c.name, c.id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Conversation) Delete() error {
	_, err := DatabaseConnection.Exec(DELETE_CONVERSATION_SQL, c.id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conversation) Id() string {
	return c.id
}

func (c *Conversation) Name() string {
	return c.name
}

func (c *Conversation) SetName(name string) {
	c.name = name
}

func (c *Conversation) CreateDate() string {
	return c.createDate
}
