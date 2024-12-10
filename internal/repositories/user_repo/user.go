package user_repo

import (
	"errors"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if record not found, which is expected behavior when the user doesn't exist
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CheckUserExistsByEmail checks if a user exists by email and returns the user if found.
func (r *UserRepository) CheckUserExistsByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if record not found, which is expected behavior when the user doesn't exist
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(user *models.User) error {
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
