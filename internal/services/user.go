package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/utils"
	"time"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(db),
	}
}

func (s *UserService) SignUp(user *models.User) (string, error) {
	// Validate National ID
	if !validateIranianNationalCode(user.NationalID) {
		return "", errors.New("invalid national ID")
	}

	// Check if user exists
	exists, err := s.repo.UserExists(user.NationalID, user.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	// Set additional fields
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save user
	if err := s.repo.SaveUser(user); err != nil {
		return "", err
	}

	// Generate JWT Token
	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
