package main

import (
	"GoServer/internal/server"
	"log"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}
