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
	rows, err := DB.Query("SELECT SYSDATE()")
	if err != nil {
		log.Fatal(err)
	}
	var date string
	for rows.Next() {
		rows.Scan(&date)
		fmt.Fprintf(w, "time: %s", date)
	}
	rows.Close()
}

func main() {
	http.HandleFunc("/view/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
