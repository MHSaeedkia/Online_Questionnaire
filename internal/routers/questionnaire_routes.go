package routers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/repositories"
)

func SetupRoutes(app *fiber.App) {
	cfg, err := config.LoadConfig("./configs/")
	if err != nil {
		log.Fatal(err)
	}

	//database connect
	DB, err := db.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(DB, "is connected successfully")

	questionnaireRepo := repositories.NewQuestionnaireRepository(DB)
	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireRepo)

	api := app.Group("/api/v1")
	questionnaireRoutes := api.Group("/questionnaires")

	questionnaireRoutes.Post("/questionnaire", questionnaireHandler.CreateQuestionnaire)
	questionnaireRoutes.Post("/question", questionnaireHandler.CreateQuestion)
	questionnaireRoutes.Post("/answer", questionnaireHandler.CreateAnswer)
}
