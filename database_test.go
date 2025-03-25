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

	user, err = getUserUsername("testUsername2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	user, err = getUserLogin("testEmail2", "testPassword2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

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
