package repositories

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type QuestionRepository interface {
	CreateQuestion(question *models.Question) error
	GetQuestionnaireByID(questionnaireID uint) (*models.Questionnaire, error)
	GetQuestionByID(questionnaireID uint, questionID uint) (*models.Question, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db}
}

func (r *questionRepository) CreateQuestion(question *models.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) GetQuestionnaireByID(questionnaireID uint) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	if err := r.db.Preload("Questions").First(&questionnaire, questionnaireID).Error; err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionRepository) GetQuestionByID(questionnaireID uint, questionID uint) (*models.Question, error) {
	var question models.Question
	err := r.db.Where("id = ? AND questionnaire_id = ?", questionID, questionnaireID).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}
