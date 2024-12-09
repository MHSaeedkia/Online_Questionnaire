package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	config "online-questionnaire/configs"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/utils"
)

// OAuthService struct
type OAuthService struct {
	ClientID     string
	ClientSecret string
	repository   *repositories.UserRepository
}

// GoogleUser contains the fields returned from the Google OAuth API
type GoogleUser struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

// NewOAuthService initializes a new OAuthService instance
func NewOAuthService(clientID, clientSecret string, repository *repositories.UserRepository) *OAuthService {
	return &OAuthService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		repository:   repository,
	}
}

// ValidateGoogleToken validates the Google OAuth token and retrieves user info
func (s *OAuthService) ValidateGoogleToken(googleToken string) (*GoogleUser, error) {
	url := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + googleToken

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error contacting Google OAuth API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token received from Google")
	}

	// Parse the response
	var userInfo GoogleUser
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &userInfo); err != nil {
		log.Printf("Error parsing Google user info: %v", err)
		return nil, err
	}

	return &userInfo, nil
}

// GetOrCreateUser checks if the user exists, if not, creates a new user
func (s *OAuthService) GetOrCreateUser(userInfo *GoogleUser) (*models.User, error) {
	// Check if the user already exists in the database by their email
	existingUser, err := s.repository.CheckUserExistsByEmail(userInfo.Email)
	if err != nil {
		return nil, err
	}

	// If user does not exist, create a new one
	if existingUser == nil {
		newUser := &models.User{
			Email:     userInfo.Email,
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Role:      models.Guest,
		}

		if err := s.repository.CreateUser(newUser); err != nil {
			return nil, err
		}

		return newUser, nil
	}

	// If the user exists, return the user
	return existingUser, nil
}

// GenerateJWTToken generates a JWT token for a given user
func (s *OAuthService) GenerateJWTToken(username, role string) (string, error) {
	tokenData, err := utils.GenerateJWTToken(username, role, config.Config{})
	if err != nil {
		return "", err
	}
	return tokenData.Token, nil
}
