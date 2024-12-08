package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *models.Questionnaire) error
	GetQuestionnaire(qustionnareId uint) (*models.Questionnaire, error)
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

func (r *questionnaireRepository) GetQuestionnaire(qustionnareId uint) (*models.Questionnaire, error) {
	var questionnare *models.Questionnaire
	return questionnare, r.db.Find(&questionnare).Where("id = ?", qustionnareId).Error
}
