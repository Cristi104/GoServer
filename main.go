package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	page, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("hello world")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
