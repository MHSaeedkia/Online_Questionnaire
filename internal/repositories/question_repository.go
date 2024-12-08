package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(question *models.Question) error
	GetQuestion(questionnaireId uint, questionId uint) (*models.Question, error)
	UpdateQuestion(quesion *models.Question) error
	DeleteQuestion(quesionId uint) error
	GetQuestionnaireByID(questionnaireID uint) (*models.Questionnaire, error)
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

func (r *questionRepository) GetQuestion(questionnaireID uint, questionID uint) (*models.Question, error) {
	var question models.Question
	err := r.db.Where("id = ? AND questionnaire_id = ?", questionID, questionnaireID).First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) UpdateQuestion(quesion *models.Question) error {
	return r.db.Updates(quesion).Where("id = ?", quesion.ID).Error
}

func (r *questionRepository) DeleteQuestion(quesionId uint) error {
	return r.db.Delete(&models.Question{}).Where("id = ?", quesionId).Error
}

func (r *questionRepository) GetQuestionnaireByID(questionnaireID uint) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	if err := r.db.Preload("Questions").First(&questionnaire, questionnaireID).Error; err != nil {
		return nil, err
	}
	return &questionnaire, nil
}
