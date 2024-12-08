package repositories

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type ConditionalLogicRepository interface {
	CreateConditionalLogic(logic *models.ConditionalLogic) error
	GetConditionalLogicByQuestionID(questionID uint) ([]models.ConditionalLogic, error)
}

type conditionalLogicRepository struct {
	db *gorm.DB
}

func NewConditionalLogicRepository(db *gorm.DB) ConditionalLogicRepository {
	return &conditionalLogicRepository{db}
}

func (r *conditionalLogicRepository) CreateConditionalLogic(logic *models.ConditionalLogic) error {
	return r.db.Create(logic).Error
}

func (r *conditionalLogicRepository) GetConditionalLogicByQuestionID(questionID uint) ([]models.ConditionalLogic, error) {
	var logic []models.ConditionalLogic
	err := r.db.Where("question_id = ?", questionID).Find(&logic).Error
	return logic, err
}
