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
	id       int64
	username string
	email    string
	password string
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
		err := rows.Scan(&user.id, &user.username, &user.email, &user.password)
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
	err := DB.QueryRow("SELECT * FROM users WHERE email = ?", email).Scan(&user.id, &user.username, &user.email, &user.password)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func makeUser(username string, email string, password string) (*User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	_, err = DB.Query("INSERT INTO users(id, username, email, password) VALUES(NULL, ?, ?, ?)", username, email, bytes)
	if err != nil {
		return nil, err
	}

	return getUserLogin(email, password)
}

func (u *User) deleteUser() error {
	_, err := DB.Query("DELETE FROM users WHERE id = ?", u.id)
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

	u.username = username

	_, err := DB.Query("UPDATE users SET username = ? WHERE id = ?", username, u.id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) setEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email given is empty")
	}

	u.email = email

	_, err := DB.Query("UPDATE users SET email = ? WHERE id = ?", email, u.id)
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

	u.password = string(bytes)

	_, err = DB.Query("UPDATE users SET password = ? WHERE id = ?", bytes, u.id)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) String() string {
	return fmt.Sprintf("%d, %s, %s, %s", u.id, u.username, u.email, u.password)
}
