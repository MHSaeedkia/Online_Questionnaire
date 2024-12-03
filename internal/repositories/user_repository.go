package repositories

import (
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Quesionnare(ownerId uint) ([]models.Questionnaire, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Quesionnare(ownerId uint) ([]models.Questionnaire, error) {
	var questionnares []models.Questionnaire
	return questionnares, r.db.Where("owner_id = ?", ownerId).Find(&questionnares).Error
}
