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
	logErr := logger.NewLogger(cfg, "questionnaire")
	if logErr != nil {
		log.Fatal(logErr)
	}

	logger.GetLogger().Info("Application started", "questionnaire", logger.Logctx{})

	// database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		message := fmt.Sprintf("Error creating database connection: %s", err.Error())
		logger.GetLogger().Fatal(message, "", logger.Logctx{}, "testtraceid")
	}

	message := fmt.Sprintf("successfully connected to database: %v", DB.Name())
	logger.GetLogger().Info(message, "", logger.Logctx{})
}
