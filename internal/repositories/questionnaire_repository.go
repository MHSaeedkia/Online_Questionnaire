package repositories

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type QuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *models.Questionnaire) error
}

type questionnaireRepository struct {
	db *gorm.DB
}

func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	return &questionnaireRepository{db}
}

func (r *questionnaireRepository) CreateQuestionnaire(questionnaire *models.Questionnaire) error {
	return r.db.Create(questionnaire).Error
}
