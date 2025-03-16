package main

import (
	"testing"
)

func TestUsers(t *testing.T) {
	user, err := NewUser("testUsername", "testEmail", "testPassword")
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

	user, err = getUserLogin("testEmail2", "testPassword2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	for i, v := range users {
		t.Logf("index: %d, %s \n", i, v)
	}

	err = user.delete()
	if err != nil {
		t.Error(err)
	}
}

func TestMessages(t *testing.T) {
	user, err := NewUser("testUsername", "testEmail", "testPassword")
	if err != nil {
		t.Error(err)
	}

	message, err := NewMessage(user, user, "test message")
	if err != nil {
		t.Error(err)
	}

	err = message.delete()
	if err != nil {
		t.Error(err)
	}

	err = user.delete()
	if err != nil {
		t.Error(err)
	}
}
