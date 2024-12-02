package services

import (
	"errors"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/utils"
)

type UserService struct {
	repository *repositories.UserRepository
	jwt        config.Config
}

func NewUserService(repository *repositories.UserRepository, cfg config.Config) *UserService {
	return &UserService{repository: repository, jwt: cfg}
}

func (s *UserService) SignUp(user *models.User) (utils.TokenData, error) {

	exist, existenceErr := s.repository.CheckUserExists(user.Email, user.NationalID)
	if existenceErr != nil {
		return utils.TokenData{}, existenceErr
	}
	if exist {
		return utils.TokenData{}, errors.New("user already exists")
	}

	isNationalCodeValid := utils.ValidateNationalID(user.NationalID)
	if !isNationalCodeValid {
		return utils.TokenData{}, errors.New("invalid national code")
	}

	user.Password = utils.GeneratePassword(user.Password)

	if createUserErr := s.repository.CreateUser(user); createUserErr != nil {
		return utils.TokenData{}, createUserErr
	}

	token, err := utils.GenerateJWTToken(user.NationalID, string(user.Role), s.jwt)
	if err != nil {
		return utils.TokenData{}, err
	}

	return token, nil
}
