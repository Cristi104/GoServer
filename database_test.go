package main

import (
	"testing"
)

func TestUsers(t *testing.T) {
	user, err := makeUser("testUsername", "testEmail", "testPassword")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	err = user.setEmail("testEmail2")
	if err != nil {
		t.Error(err)
	}
	err = user.setUsername("testUsername2")
	if err != nil {
		t.Error(err)
	}
	err = user.setPassword("testPassword2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	users, err := getUsers("testUsername2")
	if err != nil {
		t.Error(err)
	}

	for i, v := range users {
		t.Logf("index: %d, %s \n", i, v)
	}

	err = user.deleteUser()
	if err != nil {
		t.Error(err)
	}
}
