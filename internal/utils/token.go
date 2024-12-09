package utils

import (
	config "online-questionnaire/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// written by hamed maleki

func GenerateJWTToken(username, role string, cfg config.Config) (TokenData, error) {
	secret := []byte(cfg.JWT.Secret)
	expireAt := time.Now().Add(time.Minute * time.Duration(cfg.JWT.Expiration))
	claims := CustomClaims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "online-questionnaire",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return TokenData{}, err
	}

	return TokenData{
		Token:     tokenString,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: expireAt,
		Issuer:    "online-questionnaire",
	}, nil
}

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type TokenData struct {
	Token     string
	Username  string
	Role      string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Issuer    string
}
