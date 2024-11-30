package routers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/repositories"
)

func SetupRoutes(db gorm.DB, app *fiber.App) {

	questionnaireRepo := repositories.NewQuestionnaireRepository(&db)
	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireRepo)

	api := app.Group("/api")
	questionnaireRoutes := api.Group("/questionnaires")

	questionnaireRoutes.Post("/", questionnaireHandler.CreateQuestionnaire)
}
