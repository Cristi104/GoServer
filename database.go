package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

type config struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DBName string
}

func init() {
	var temp config

	// reading database connection data
	data, err := os.ReadFile("config/DB.json")
	if err != nil {
		// if the file is missing create an empty config file
		log.Printf("Missing configuration file created, please complete database connection information inside config/DB.json")

		data, err1 := json.Marshal(temp)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.WriteFile("config/DB.json", data, 0600)
		if err1 != nil {
			log.Fatal(err1)
		}

		log.Fatal(err)
	}

	err = json.Unmarshal(data, &temp)
	if err != nil {
		log.Fatal(err)
	}

	cfg := mysql.Config{
		User:   temp.User,
		Passwd: temp.Passwd,
		Net:    temp.Net,
		Addr:   temp.Addr,
		DBName: temp.DBName,
	}

	if err != nil {
		log.Fatal(err)
	}

	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to database")
}

type User struct {
	Id       int64
	Username string
	Email    string
	Password string
}

func GetUser(id int64) (*User, error) {
	var user User

	err := DB.QueryRow(`
		SELECT * 
		FROM users 
		WHERE id = ?
	`, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserUsername(username string) (*User, error) {
	var user User

	err := DB.QueryRow(`
		SELECT * 
		FROM users 
		WHERE username = ?
	`, username).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserLogin(email string, password string) (*User, error) {
	var user User

	err := DB.QueryRow(`
		SELECT * 
		FROM users 
		WHERE email = ?
	`, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewUser(username string, email string, password string) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	result, err := DB.Exec(`
		INSERT INTO users(id, username, email, password) 
		VALUES(NULL, ?, ?, ?)
	`, username, email, bytes)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetUser(id)
}

func (u *User) Delete() error {
	_, err := DB.Query(`
		DELETE FROM users 
		WHERE id = ?
	`, u.Id)
	if err != nil {
		return err
	}

	*u = User{}

	return nil
}

func (u *User) SetUsername(username string) error {
	if len(username) == 0 {
		return errors.New("username given is empty")
	}

	u.Username = username

	_, err := DB.Exec(`
		UPDATE users 
		SET username = ? 
		WHERE id = ?
	`, username, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SetEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email given is empty")
	}

	u.Email = email

	_, err := DB.Exec(`
		UPDATE users 
		SET email = ? 
		WHERE id = ?
	`, email, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password given is empty")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	_, err = DB.Exec(`
		UPDATE users 
		SET password = ? 
		WHERE id = ?
	`, bytes, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) String() string {
	return fmt.Sprintf("%d, %s, %s, %s", u.Id, u.Username, u.Email, u.Password)
}

type ConversationMember struct {
	User
}

type Conversation struct {
	Id         int64
	Name       string
	CreateDate string
	Members    []*ConversationMember
}

func GetConversation(id int64) (*Conversation, error) {
	var conversation Conversation

	// get conversation details
	err := DB.QueryRow(`
		SELECT * 
		FROM conversations
		WHERE id = ?
	`, id).Scan(&conversation.Id, &conversation.Name, &conversation.CreateDate)
	if err != nil {
		return nil, err
	}

	// get all users in the conversation
	rows, err := DB.Query(`
		SELECT u.*
		FROM in_conversation ic
		LEFT JOIN users u ON u.id = ic.user_id
		WHERE ic.conversation_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var member ConversationMember
	for rows.Next() {
		err := rows.Scan(&member.Id, &member.Username, &member.Email, &member.Password)
		if err != nil {
			return nil, err
		}

		copy := member
		conversation.Members = append(conversation.Members, &copy)
	}

	return &conversation, nil
}

func NewConversation(creator *User, name string) (*Conversation, error) {
	// create conversation
	result, err := DB.Exec(`
		INSERT INTO conversations(id, name, create_date)
		VALUES(NULL, ?, SYSDATE())
	`)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// add creator to conversation
	result, err = DB.Exec(`
		INSERT INTO in_conversation(user_id, conversationd_id)
		VALUES(?, ?)
	`, creator.Id, id)
	if err != nil {
		return nil, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetConversation(id)
}

func GetUserConversations(user *User) ([]*Conversation, error) {
	rows, err := DB.Query(`
		SELECT conversation_id
		FROM in_conversation
		WHERE user_id = ?
	`, user.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*Conversation
	var id int64

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		conversation, err := GetConversation(id)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func UserInConversation(userId int64, conversationId int64) bool {
	var result int64
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM in_conversation
		WHERE user_id = ?
		AND conversation_id = ?
	`, userId, conversationId).Scan(&result)
	if err != nil {
		log.Println(err)
		return false
	}

	return result != 0
}

type Message struct {
	Id             int64
	SendDate       string
	Body           string
	SenderId       int64
	ConversationId int64
}

func GetMessage(id int64) (*Message, error) {
	var message Message

	err := DB.QueryRow(`
		SELECT * 
		FROM messages 
		WHERE id = ?
	`, id).Scan(&message.Id, &message.SendDate, &message.Body, &message.SenderId, &message.ConversationId)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func NewMessage(senderId int64, conversationId int64, body string) (*Message, error) {
	result, err := DB.Exec(`
		INSERT INTO messages(id, send_date, body, sender_id, conversation_id) 
		VALUES(NULL, SYSDATE(), ?, ?, ?)
	`, body, senderId, conversationId)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetMessage(id)
}

func GetConversationMessages(conversationId int64) ([]*Message, error) {
	rows, err := DB.Query(`
		SELECT *
		FROM messages
		WHERE conversation_id = ?
		ORDER BY send_date
	`, conversationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	var message Message

	for rows.Next() {
		err = rows.Scan(&message.Id, &message.SendDate, &message.Body, &message.SenderId, &message.ConversationId)
		if err != nil {
			return nil, err
		}

		copy := message
		messages = append(messages, &copy)
	}

	return messages, nil
}

func (m *Message) Delete() error {
	_, err := DB.Exec(`
		DELETE FROM messages 
		WHERE id = ?
	`, m.Id)
	if err != nil {
		return err
	}

	return nil
}
