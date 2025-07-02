package tests

import (
	"testing"

	"GoServer/repository"
)

func TestConnection(t *testing.T) {
	var ret string

	err := repository.DatabaseConnection.QueryRow("SELECT CURRENT_TIMESTAMP").Scan(&ret)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUser(t *testing.T) {
	user, err := repository.InsertUser("test_username", "test@email.com", "test_pass")
	if err != nil {
		t.Fatal(err)
	}

	user.Nickname = "test_nick"

	err = user.Update()
	if err != nil {
		t.Fatal(err)
	}

	user, err = repository.SelectUserById(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = user.Delete()
	if err != nil {
		t.Fatal(err)
	}

}
