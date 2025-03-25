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

	err = user.SetEmail("testEmail2")
	if err != nil {
		t.Error(err)
	}
	err = user.SetUsername("testUsername2")
	if err != nil {
		t.Error(err)
	}
	err = user.SetPassword("testPassword2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	user, err = GetUserUsername("testUsername2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	user, err = GetUserLogin("testEmail2", "testPassword2")
	if err != nil {
		t.Error(err)
	}
	t.Log(user)

	err = user.Delete()
	if err != nil {
		t.Error(err)
	}
}

func TestMessages(t *testing.T) {
	// user, err := NewUser("testUsername", "testEmail", "testPassword")
	// if err != nil {
	// 	t.Error(err)
	// }

	// message, err := NewMessage(user, user, "test message")
	// if err != nil {
	// 	t.Error(err)
	// }

	// err = message.delete()
	// if err != nil {
	// 	t.Error(err)
	// }

	// err = user.Delete()
	// if err != nil {
	// 	t.Error(err)
	// }
}
