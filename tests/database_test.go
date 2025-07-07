package tests

import (
	"fmt"
	"slices"
	"strings"
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

	defer func() {
		err = user.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	user.Nickname = "test_nickname"

	err = user.Update()
	if err != nil {
		user.Delete()
		t.Fatal(err)
	}

	user, err = repository.SelectUserById(user.Id)
	if err != nil {
		user.Delete()
		t.Fatal(err)
	}
}

func TestConversation(t *testing.T) {
	// creating 3 users
	// user 1
	user1, err := repository.InsertUser("test_username1", "test1@email.com", "test_pass1")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = user1.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// user 2
	user2, err := repository.InsertUser("test_username2", "test2@email.com", "test_pass2")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = user2.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// user 3
	user3, err := repository.InsertUser("test_username3", "test3@email.com", "test_pass3")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = user3.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// create a conversation
	conv, err := repository.InsertConversation("conversation1")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = conv.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// add multiple users to conversation
	err = conv.AddUsers([]string{user1.Id, user2.Id})
	if err != nil {
		t.Fatal(err)
	}

	// test user 1 in conversation
	convs, err := repository.SelectConversationsByUser(user1.Id)
	fmt.Println(convs)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(convs, func(c repository.Conversation) bool {
		return strings.Compare(c.Id, conv.Id) == 0
	}) {
		t.Fatal("expected user1 to be in conversation1")
	}
	if len(convs) != 1 {
		t.Fatal("expected user1 to be only in conversation1")
	}

	// test user 3 not in conversation
	convs, err = repository.SelectConversationsByUser(user3.Id)
	if err != nil {
		t.Fatal(err)
	}

	if slices.ContainsFunc(convs, func(c repository.Conversation) bool {
		return strings.Compare(c.Id, conv.Id) == 0
	}) {
		t.Fatal("expected user3 to not be in conversation1")
	}
	if len(convs) != 0 {
		t.Fatal("expected user1 to not be in any conversation")
	}

	// test conversation update
	conv.Name = "new name"
	err = conv.Update()
	if err != nil {
		t.Fatal(err)
	}

	// add user 3 to conversation
	err = conv.AddUser(user3.Id)
	if err != nil {
		t.Fatal(err)
	}

	// test that he is actualy added
	convs, err = repository.SelectConversationsByUser(user3.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(convs, func(c repository.Conversation) bool {
		return strings.Compare(c.Id, conv.Id) == 0
	}) {
		t.Fatal("expected user3 to be in conversation1")
	}
	if len(convs) != 1 {
		t.Fatal("expected user3 to be only in conversation1")
	}
}

func TestMessage(t *testing.T)  {
	// creating 2 users
	// user 1
	user1, err := repository.InsertUser("test_username1", "test1@email.com", "test_pass1")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = user1.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// user 2
	user2, err := repository.InsertUser("test_username2", "test2@email.com", "test_pass2")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = user2.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// create a conversation
	conv, err := repository.InsertConversation("conversation1")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = conv.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// add users to conversation
	err = conv.AddUsers([]string{user1.Id, user2.Id})
	if err != nil {
		t.Fatal(err)
	}

	// send 2 messages
	message1, err := repository.InsertMessage("test message from user 1", user1.Id, conv.Id)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = message1.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	message2, err := repository.InsertMessage("test message from user 2", user2.Id, conv.Id)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = message2.Delete()
		if err != nil {
			t.Fatal(err)
		}
	}()

	message2.Body = "i have edited this message"
	err = message2.Update()
	if err != nil {
		t.Fatal(err)
	}

	messages, err := repository.SelectMessagesByConversation(conv.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.ContainsFunc(messages, func(m repository.Message) bool {
		return m.Id == message1.Id
	}) {
		t.Fatal("message 1 not in conversation")
	}
	if !slices.ContainsFunc(messages, func(m repository.Message) bool {
		return m.Id == message2.Id
	}) {
		t.Fatal("message 2 not in conversation")
	} else {
		if !slices.ContainsFunc(messages, func(m repository.Message) bool{
			return m.Body == message2.Body
		}) {
			t.Fatal("message 2 was not edited")
		}
	}
}




