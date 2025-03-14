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

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/home.html")
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlerFunc)
	http.HandleFunc("/view/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
