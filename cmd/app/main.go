package main

import (
	"fmt"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig("../../configs/")
	if err != nil {
		log.Fatal(err)
	}
	log := logger.NewLogger(cfg, "", "")
	// database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		message := fmt.Sprintf("Error creating database connection: %s", err.Error())
		log.Fatal(message, "", nil, "testtraceid")
	}

	message := fmt.Sprintf("successfully connected to database: %v", DB.Name())
	log.Info(message, "", nil)
}
