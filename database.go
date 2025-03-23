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

func getUser(id int64) (*User, error) {
	var user User

	err := DB.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getUsers(username string) ([]*User, error) {
	rows, err := DB.Query("SELECT * FROM users WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		copy := user
		users = append(users, &copy)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func getUserLogin(email string, password string) (*User, error) {
	var user User
	err := DB.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
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

	result, err := DB.Exec("INSERT INTO users(id, username, email, password) VALUES(NULL, ?, ?, ?)", username, email, bytes)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return getUser(id)
}

func getFriends(user *User) ([]*User, error) {
	rows, err := DB.Query(`
	SELECT * 
	FROM users 
	WHERE id IN (
		SELECT 
			CASE 
				WHEN receiver_id = ? THEN sender_id 
				ELSE receiver_id 
			END 
		FROM messages 
		WHERE sender_id = ? 
		OR receiver_id = ?
	)`, user.Id, user.Id, user.Id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []*User

	var tempUser User
	for rows.Next() {
		err := rows.Scan(&tempUser.Id, &tempUser.Username, &tempUser.Email, &tempUser.Password)
		if err != nil {
			return nil, err
		}
		copy := tempUser
		users = append(users, &copy)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (u *User) delete() error {
	_, err := DB.Query("DELETE FROM users WHERE id = ?", u.Id)
	if err != nil {
		return err
	}

	*u = User{}

	return nil
}

func (u *User) setUsername(username string) error {
	if len(username) == 0 {
		return errors.New("username given is empty")
	}

	u.Username = username

	_, err := DB.Exec("UPDATE users SET username = ? WHERE id = ?", username, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) setEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email given is empty")
	}

	u.Email = email

	_, err := DB.Exec("UPDATE users SET email = ? WHERE id = ?", email, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password given is empty")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	u.Password = string(bytes)

	_, err = DB.Exec("UPDATE users SET password = ? WHERE id = ?", bytes, u.Id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) String() string {
	return fmt.Sprintf("%d, %s, %s, %s", u.Id, u.Username, u.Email, u.Password)
}

type Message struct {
	Id         int64
	SendDate   string
	Body       string
	SenderId   int64
	ReceiverId int64
}

func getMessage(id int64) (*Message, error) {
	var message Message

	err := DB.QueryRow("SELECT * FROM messages WHERE id = ?", id).Scan(&message.Id, &message.SendDate, &message.Body, &message.SenderId, &message.ReceiverId)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func NewMessage(sender *User, receiver *User, body string) (*Message, error) {
	result, err := DB.Exec("INSERT INTO messages(id, send_date, body, sender_id, receiver_id) VALUES(NULL, SYSDATE(), ?, ?, ?)", body, sender.Id, receiver.Id)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return getMessage(id)
}

func (m *Message) delete() error {
	_, err := DB.Exec("DELETE FROM messages WHERE id = ?", m.Id)
	if err != nil {
		return err
	}

	return nil
}
