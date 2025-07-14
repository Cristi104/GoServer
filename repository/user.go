package repository

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string
	Username     string
	Nickname     string
	Email        string
	PasswordHash string
}

const INSERT_USER_SQL = "INSERT INTO users(username, nickname, email, password) VALUES($1, $1, $2, $3) RETURNING id;"
const INSERT_FRIEND_SQL = "INSERT INTO is_friend(user1_id, user2_id) VALUES($1, $2);"
const SELECT_USER_BY_ID_SQL = "SELECT * FROM users WHERE id = $1;"
const SELECT_USER_BY_SIGN_IN_SQL = "SELECT * FROM users WHERE email = $1 OR username = $1;"
const SELECT_ALL_USERS_IN_CONVERSATION_SQL = "SELECT u.* FROM in_conversation c LEFT JOIN users u ON u.id = c.user_id WHERE c.conversation_id = $1;"
const SELECT_ALL_USER_FRIENDS_SQL = `
SELECT u.* 
FROM users u 
WHERE u.id IN (
	SELECT CASE
		WHEN if.user1_id = $1 THEN if.user2_id
		WHEN if.user2_id = $1 THEN if.user1_id
	END
	FROM is_friend if 
	WHERE if.user1_id = $1 OR if.user2_id = $1
);`
const SELECT_ALL_USERS_BY_USERNAME_SQL = "SELECT * FROM users WHERE username LIKE $1;"
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

	err := DatabaseConnection.QueryRow(SELECT_USER_BY_ID_SQL, id).Scan(&user.Id, &user.Username, &user.Nickname, &user.Email, &user.PasswordHash)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func SelectUserBySignIn(email, password string) (User, error) {
	var user User

	err := DatabaseConnection.QueryRow(SELECT_USER_BY_SIGN_IN_SQL, email).Scan(&user.Id, &user.Username, &user.Nickname, &user.Email, &user.PasswordHash)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func SelectUserFriends(userId string) ([]User, error) {
	rows, err := DatabaseConnection.Query(SELECT_ALL_USER_FRIENDS_SQL, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Nickname, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
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
		err := rows.Scan(&user.Id, &user.Username, &user.Nickname, &user.Email, &user.PasswordHash)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func SelectUsersByUsername(username string) ([]User, error) {
	rows, err := DatabaseConnection.Query(SELECT_ALL_USERS_BY_USERNAME_SQL, fmt.Sprintf("%s%s%s", "%", username, "%"))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Nickname, &user.Email, &user.PasswordHash)
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

func (u *User) AddFriend(friendId string) error {
	_, err := DatabaseConnection.Exec(INSERT_FRIEND_SQL, u.Id, friendId)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Update() error {
	_, err := DatabaseConnection.Exec(UPDATE_USER_SQL, u.Username, u.Nickname, u.Email, u.PasswordHash, u.Id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	_, err := DatabaseConnection.Exec(DELETE_USER_SQL, u.Id)
	if err != nil {
		return err
	}

	*u = User{}
	return nil
}

func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	u.PasswordHash = hash
	return nil
}
