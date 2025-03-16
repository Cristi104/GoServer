package main

import (
	"fmt"
	"log"
)

func testDatabase() {
	testUsers()
}

func testUsers() {
	user, err := makeUser("testUsername", "testEmail", "testPassword")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	err = user.setEmail("testEmail2")
	if err != nil {
		log.Fatal(err)
	}
	err = user.setUsername("testUsername2")
	if err != nil {
		log.Fatal(err)
	}
	err = user.setPassword("testPassword2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	users, err := getUsers("testUsername2")
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range users {
		fmt.Printf("index: %d, %s \n", i, v)
	}

	err = user.deleteUser()
	if err != nil {
		log.Fatal(err)
	}
}
