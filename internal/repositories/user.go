package repositories

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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
	if err := r.db.Where("national_id = ?", nationalID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if user is not found
		}
		return nil, err // Return error if any other issue occurs
	}
	return &user, nil // Return the user object if found
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// VerifyPassword compares the stored hash with the provided password.
func (r *UserRepository) VerifyPassword(storedPassword, inputPassword string) bool {
	// Assuming the stored password is a hash (bcrypt)
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(inputPassword))
	return err == nil
}
