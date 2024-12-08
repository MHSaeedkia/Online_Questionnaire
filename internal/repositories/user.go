package repositories

import (
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CheckUserExists checks if a user exists by national ID and returns the user if found.
func (r *UserRepository) CheckUserExists(nationalID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("national_id = ?", nationalID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	// Ensure that only the hashed password is sent to the database
	if err := r.db.Create(&models.User{
		NationalID:    user.NationalID,
		Email:         user.Email,
		Password:      user.Password,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Gender:        user.Gender,
		DateOfBirth:   user.DateOfBirth,
		City:          user.City,
		WalletBalance: 0,
		Role:          models.Guest,
	}).Error; err != nil {
		return err
	}
	return nil
}
