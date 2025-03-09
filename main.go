package main

import (
	"fmt"
)

func main() {
	page1 := &Page{Title: "Test", Body: []byte("Lorem ipsum.")}
	page1.save()
	page2, _ := loadPage("Test")
	fmt.Println(string(page2.Body))
	fmt.Println("hello world")
}
