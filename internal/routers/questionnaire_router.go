package routers

import (
	"online-questionnaire/internal/handlers"
	"online-questionnaire/internal/repositories"

	"gorm.io/gorm"
)

func questionnareRouter(db *gorm.DB) *handlers.QuestionnaireHandler {
	questionnaireRepo := repositories.NewQuestionnaireRepository(db)
	return handlers.NewQuestionnaireHandler(questionnaireRepo)
}
