package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func saveAsJson(a interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func loadFromJson(a interface{}, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, a)
}

type config struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DBName string
}

func init() {
	var temp config
	fmt.Printf("%s\n%s\n%s\n%s\n%s\n", temp.User, temp.Passwd, temp.Net, temp.Addr, temp.DBName)
	err := loadFromJson(&temp, "config.json")

	if err != nil {
		log.Printf("Missing configuration file created, please complete database connection information inside config.json")
		saveAsJson(temp, "config.json")
		log.Fatal(err)
	}

	cfg := mysql.Config{
		User:   temp.User,
		Passwd: temp.Passwd,
		Net:    temp.Net,
		Addr:   temp.Addr,
		DBName: temp.DBName,
	}

	if err != nil {
		log.Fatal(err)
	}

	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected to database")
}
