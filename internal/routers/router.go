package routers

import (
	"fmt"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/db"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// get config .
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

	api := app.Group("/api/v1")
	questionnaireHandler := questionnareRouter(DB)
	userRepository := userRouter(DB)

	questionnaireRoutes := api.Group("/questionnaires")
	questionnaireRoutes.Post("/questionnaire", questionnaireHandler.CreateQuestionnaire)
	questionnaireRoutes.Post("/question", questionnaireHandler.CreateQuestion)
	questionnaireRoutes.Post("/answer", questionnaireHandler.CreateAnswer)

	userRouter := api.Group("/user")
	userRouter.Get("/questionnaires", userRepository.Quesionnare)
	userRouter.Put("/questionnaires/edit", userRepository.EditQuestionnare)

}
