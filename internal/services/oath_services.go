package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"online-questionnaire/configs"
	"online-questionnaire/internal/models"
	"online-questionnaire/internal/repositories"
	"online-questionnaire/internal/utils"
)

// OAuthService struct
type OAuthService struct {
	ClientID     string
	ClientSecret string
	Config       config.Config
	repository   *repositories.UserRepository
}

// NewOAuthService initializes a new OAuthService instance
func NewOAuthService(cfg config.Config, repository *repositories.UserRepository) *OAuthService {
	return &OAuthService{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Config:       cfg,
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

	// Check if the user already exists in the database by their email
	existingUser, err := s.repository.CheckUserExistsByEmail(userInfo.Email)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		// If the user does not exist, create a new one
		user := &models.User{
			Email:     userInfo.Email,
			FirstName: userInfo.GivenName,
			LastName:  userInfo.FamilyName,
			Role:      models.Guest,
		}
		if err := s.repository.CreateUser(user); err != nil {
			return nil, err
		}
	}

	return &userInfo, nil
}

// GoogleUser contains the fields returned from the Google OAuth API
type GoogleUser struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

// GenerateJWTToken generates a JWT token for a given user
func (s *OAuthService) GenerateJWTToken(username, role string) (string, error) {
	tokenData, err := utils.GenerateJWTToken(username, role, config.Config{})
	if err != nil {
		return "", err
	}
	return tokenData.Token, nil
}

// Authenticate generates an access token for OAuth authentication
func (s *OAuthService) Authenticate(clientID, clientSecret string) (utils.TokenData, error) {
	if clientID != s.ClientID || clientSecret != s.ClientSecret {
		return utils.TokenData{}, errors.New("invalid client credentials")
	}

	tokenData, err := utils.GenerateJWTToken(clientID, "oauth-client", s.Config)
	if err != nil {
		return utils.TokenData{}, err
	}

	return tokenData, nil
}

// ValidateToken uses the utility function to validate a token
func (s *OAuthService) ValidateToken(tokenString string) (*utils.CustomClaims, error) {
	return utils.ValidateToken(tokenString, s.Config.JWT.Secret)
}
