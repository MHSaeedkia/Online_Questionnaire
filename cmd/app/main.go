package main

import (
	"fmt"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/jwtgen"
)

func main() {
	cfg, err := config.LoadConfig("../../configs/")
	if err != nil {
		log.Fatal(err)
	}
	token, err := jwtgen.GenerateJWTToken("hamed", cfg)
	fmt.Println(token)
	//database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")
}
