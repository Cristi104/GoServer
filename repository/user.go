package repository

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id           string
	username     string
	nickname     string
	email        string
	passwordHash string
}

const INSERT_USER_SQL = "INSERT INTO users(username, nickname, email, password) VALUES($1, $1, $2, $3) RETURNING id;"
const SELECT_USER_BY_ID_SQL = "SELECT * FROM users WHERE id = $1;"
const SELECT_ALL_USERS_IN_CONVERSATION_SQL = "SELECT u.* FROM in_conversation c LEFT JOIN users u ON u.id = c.user_id;"
const UPDATE_USER_SQL = "UPDATE users SET username = $1, nickname = $2, email = $3, password = $4 WHERE id = $5;"
const DELETE_USER_SQL = "DELETE FROM users WHERE id = $1;"

func InsertUser(username string, email string, password string) (User, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	var id string
	err = DatabaseConnection.QueryRow(INSERT_USER_SQL, username, email, string(bytes)).Scan(&id)
	if err != nil {
		return User{}, err
	}

	return User{id, username, username, email, string(bytes)}, nil
}

func SelectUserById(id string) (User, error) {
	var user User

	err := DatabaseConnection.QueryRow(SELECT_USER_BY_ID_SQL, id).Scan(&user.id, &user.username, &user.nickname, &user.email, &user.passwordHash)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func SelectUsersInConversation(conversationId string) ([]User, error) {
	rows, err := DatabaseConnection.Query(SELECT_ALL_USERS_IN_CONVERSATION_SQL, conversationId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	var user User
	for rows.Next() {
		err := rows.Scan(&user.id, &user.username, &user.email, &user.passwordHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (u *User) Update() error {
	_, err := DatabaseConnection.Exec(UPDATE_USER_SQL, u.username, u.nickname, u.email, u.passwordHash, u.id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	_, err := DatabaseConnection.Exec(DELETE_USER_SQL, u.id)
	if err != nil {
		return err
	}

	*u = User{}
	return nil
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Nickname() string {
	return u.nickname
}

func (u *User) Email() string {
	return u.email
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) SetUsername(username string) {
	u.username = username
}

func (u *User) SetNickname(nickname string) {
	u.nickname = nickname
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	u.passwordHash = hash;
	return nil
}

