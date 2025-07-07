package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
)

var DatabaseConnection *sql.DB

type config struct {
	Host   string
	Port   int32
	User   string
	Passwd string
	DBName string
}

func init() {
	temp := config{"198.7.0.2", 5432, "root", "root_pass", "GoServer"}

	data, err := os.ReadFile("config/database_connection.json")
	if err != nil {
		data, err1 := json.Marshal(temp)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.MkdirAll(filepath.Join(".", "config"), os.ModePerm)
		if err1 != nil {
			log.Fatal(err1)
		}

		err1 = os.WriteFile("config/database_connection.json", data, 0600)
		if err1 != nil {
			log.Fatal(err1)
		}

		log.Println("Missing configuration file created, if the database connection infromation was changed from the dafault update it inside config/database_connection.json")
	}

	err = json.Unmarshal(data, &temp)
	if err != nil {
		log.Fatal(err)
	}

	data_string := fmt.Sprintf(
		"host=%s "+
			"port=%d "+
			"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"sslmode=disable",
		temp.Host, temp.Port, temp.User, temp.Passwd, temp.DBName)

	DatabaseConnection, err = sql.Open("postgres", data_string)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to database: %s at %s:%d\n", temp.DBName, temp.Host, temp.Port)
}
