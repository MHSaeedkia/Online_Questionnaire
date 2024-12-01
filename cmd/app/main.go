package main

import (
	"fmt"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
)

func main() {
	cfg, err := config.LoadConfig("../../configs/")
	if err != nil {
		log.Fatal(err)
	}
	//database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")
}
