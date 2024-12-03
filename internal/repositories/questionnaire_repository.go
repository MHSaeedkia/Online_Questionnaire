package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *models.Questionnaire) error
	CreateQuestion(quesion *models.Question) error
	CreateAnswer(response *models.Response) error
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

func (r *questionnaireRepository) CreateQuestion(quesion *models.Question) error {
	return r.db.Create(quesion).Error
}

func (r *questionnaireRepository) CreateAnswer(response *models.Response) error {
	return r.db.Create(response).Error
}
