package questionnaire_repo

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type OptionRepository interface {
	CreateOptions(options []models.Option) error
	GetOptionByID(optionID uint) (*models.Option, error)
}

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{db}
}

func (r *optionRepository) GetOptionByID(optionID uint) (*models.Option, error) {
	var option models.Option
	err := r.db.First(&option, optionID).Error
	if err != nil {
		return nil, err
	}
	return &option, nil
}

func (r *optionRepository) CreateOptions(options []models.Option) error {
	return r.db.Create(&options).Error
}
