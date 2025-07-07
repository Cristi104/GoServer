package handler

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var JWTKey = ""

func init() {
	type jwtJson struct {
		Key string
	}

	token := jwtJson{"default_token"}

	readData, err := os.ReadFile("config/jwt.json")
	if err != nil {
		data, err1 := json.Marshal(token)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.MkdirAll(filepath.Join(".", "config"), os.ModePerm)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.WriteFile("config/jwt.json", data, 0600)
		if err1 != nil {
			log.Fatal(err1)
		}
		log.Println("Missing configuration file created, please change the JWT key information inside config/jwt.json")
		
		JWTKey = token.Key
		return
	}

	err = json.Unmarshal(readData, &token)
	if err != nil {
		log.Fatal(err)
	}

	if token.Key == "default_token" {
		log.Println("Warning: JWT key is set to the default value, please change it.")
	}

	JWTKey = token.Key


}
