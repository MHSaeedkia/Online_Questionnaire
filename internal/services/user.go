package services

import (
	"errors"
	"log"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/utils"
	"time"
)

type UserService struct {
	repository *repositories.UserRepository
	jwt        config.Config
}

func NewUserService(repository *repositories.UserRepository, cfg config.Config) *UserService {
	return &UserService{repository: repository, jwt: cfg}
}

func (s *UserService) SignUp(user *models.User) (utils.TokenData, error) {
	// Check if the user already exists
	exist, existenceErr := s.repository.CheckUserExists(user.NationalID)
	if existenceErr != nil {
		return utils.TokenData{}, existenceErr
	}
	if exist != nil {
		return utils.TokenData{}, errors.New("user already exists")
	}

	// Validate the national code
	isNationalCodeValid := utils.ValidateNationalID(user.NationalID)
	if !isNationalCodeValid {
		return utils.TokenData{}, errors.New("invalid national code")
	}

	// Hash the user's password
	hashedPassword, err := utils.GeneratePassword(user.Password)
	if err != nil {
		return utils.TokenData{}, err
	}
	user.Password = hashedPassword
	// Process the date of birth
	loc, _ := time.LoadLocation("Asia/Tehran") // Set to your desired location

	// Check if the DateOfBirth is in the correct format and append time
	formattedDate := user.DateOfBirth.Format("2006-01-02") + "T00:00:00Z" // Append T00:00:00Z
	log.Println("Formatted DateOfBirth:", formattedDate)
	parsedDate, err := time.Parse(time.RFC3339, formattedDate) // RFC3339 covers full datetime format
	if err != nil {
		return utils.TokenData{}, errors.New("invalid date of birth format, expected YYYY-MM-DD")
	}

	// Set the parsed date with timezone
	user.DateOfBirth = parsedDate.In(loc)

	// Create the user in the database
	if createUserErr := s.repository.CreateUser(user); createUserErr != nil {
		return utils.TokenData{}, createUserErr
	}

	// Generate JWT token
	token, err := utils.GenerateJWTToken(user.NationalID, string(user.Role), s.jwt)
	if err != nil {
		return utils.TokenData{}, err
	}

	return token, nil
}

func (s *UserService) Login(nationalID, password string) (utils.TokenData, error) {
	// Fetch the user by national ID
	user, err := s.repository.CheckUserExists(nationalID)
	if err != nil {
		return utils.TokenData{}, errors.New("user not found")
	}

	// Compare the password
	if !utils.ComparePassword(user.Password, password) {
		return utils.TokenData{}, errors.New("incorrect password")
	}

	// Generate the JWT token
	token, err := utils.GenerateJWTToken(user.NationalID, string(user.Role), s.jwt)
	if err != nil {
		return utils.TokenData{}, err
	}

	return token, nil
}
