package questionnaire_repo

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type QuestionnaireRepository interface {
	CreateQuestionnaire(questionnaire *models.Questionnaire) error
	GetByID(id uint) (*models.Questionnaire, error)
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

func (r *questionnaireRepository) GetByID(id uint) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	err := r.db.First(&questionnaire, id).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}
