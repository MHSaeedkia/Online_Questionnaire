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

	api := app.Group("/api")
	questionnaireRoutes := api.Group("/questionnaires")

	questionnaireRepo := repositories.NewQuestionnaireRepository(DB)
	questionRepo := repositories.NewQuestionRepository(DB)
	optionRepo := repositories.NewOptionRepository(DB)
	conditionalLogicRepo := repositories.NewConditionalLogicRepository(DB)

	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireRepo)
	questionHandler := handlers.NewQuestionHandler(questionRepo)
	optionHandler := handlers.NewOptionHandler(optionRepo, questionRepo)
	conditionalLogicHandler := handlers.NewConditionalLogicHandler(conditionalLogicRepo, questionRepo, optionRepo)

	questionnaireRoutes.Post("/", questionnaireHandler.CreateQuestionnaire)

	questionnaireRoutes.Post("/:questionnaire_id/questions", questionHandler.CreateQuestion)

	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/options", optionHandler.CreateOptions)

	questionnaireRoutes.Post("/:questionnaire_id/questions/:question_id/conditional-logic", conditionalLogicHandler.CreateConditionalLogic)
}
