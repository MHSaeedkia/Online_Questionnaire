package services

import (
	"errors"
	"fmt"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/logger"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories/user_repo"
	"online-questionnaire/internal/utils"
)

type UserService struct {
	repository *user_repo.UserRepository
	jwt        config.Config
}

func NewUserService(repository *user_repo.UserRepository, cfg config.Config) *UserService {
	return &UserService{repository: repository, jwt: cfg}
}

func (s *UserService) SignUp(user *models.User) (utils.TokenData, error) {
	// Check if the user already exists
	exist, existenceErr := s.repository.CheckUserExists(user.NationalID)
	if existenceErr != nil {
		return utils.TokenData{}, existenceErr
	}
	if exist != nil {
		logger.GetLogger().Info(fmt.Sprintf("User already exists: %v", user.NationalID), nil, logger.Logctx{})

		return utils.TokenData{}, errors.New("user already exists")
	}

	// Validate the national code
	isNationalCodeValid := utils.ValidateNationalID(user.NationalID)
	if !isNationalCodeValid {
		logger.GetLogger().Error(fmt.Sprintf("Invalid national ID: %v", user.NationalID), nil, logger.Logctx{}, "")
		return utils.TokenData{}, errors.New("invalid national code")

	}

	// Hash the user's password
	hashedPassword, err := utils.GeneratePassword(user.Password)
	if err != nil {
		return utils.TokenData{}, err
	}
	user.Password = hashedPassword

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
		logger.GetLogger().Error(fmt.Sprintf("Incorrect username or password: %v", user.NationalID), err, logger.Logctx{}, "")
		return utils.TokenData{}, errors.New("Incorrect Username or Password")
	}

	// Compare the password
	if !utils.ComparePassword(user.Password, password) {
		logger.GetLogger().Error(fmt.Sprintf("Incorrect username or password: %v", user.NationalID), err, logger.Logctx{}, "")
		return utils.TokenData{}, errors.New("incorrect username or password")
	}

	// Generate the JWT token
	token, err := utils.GenerateJWTToken(user.NationalID, string(user.Role), s.jwt)
	if err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Couldn't get token: %v", user.NationalID), err, logger.Logctx{}, "")
		return utils.TokenData{}, err
	}
	logger.GetLogger().Info(fmt.Sprintf("Token has been generated for user: %v", user.NationalID), err, logger.Logctx{})
	return token, nil
}
