package repositories

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UserExists(nationalID, email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("national_id = ? OR email = ?", nationalID, email).Count(&count).Error
	return count > 0, err
}
func (r *UserRepository) SaveUser(user *models.User) error {
	return r.db.Create(user).Error
}
