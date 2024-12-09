package services

import (
	"fmt"
	"log"
	"math/rand"
	"online-questionnaire/internal/utils"
	"time"
)

// VerificationService struct
type VerificationService struct {
	redisClient *utils.RedisClient
	emailSender *utils.SendEmail
}

// NewVerificationService creates a new instance of VerificationService
func NewVerificationService(redisClient *utils.RedisClient, emailSender *utils.SendEmail) *VerificationService {
	return &VerificationService{
		redisClient: redisClient,
		emailSender: emailSender,
	}
}

// GenerateVerificationCode generates a random 6-digit code
func (s *VerificationService) GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

// SendVerificationCode sends a verification code to the user's email
func (s *VerificationService) SendVerificationCode(email string) error {
	// Generate a verification code
	code := s.GenerateVerificationCode()

	// Save the code to Redis with a 10-minute expiration time
	err := s.redisClient.Set(email, code, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to store verification code: %v", err)
	}

	// Send the code to the user's email
	subject := "Verification Code"
	body := fmt.Sprintf("Your verification code is: %s", code)
	err = s.emailSender.SendEmail(email, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Verification code sent to: %s", email)
	return nil
}

// ValidateCode validates the code entered by the user
func (s *VerificationService) ValidateCode(email, code string) (bool, error) {
	// Retrieve the saved code from Redis
	savedCode, err := s.redisClient.Get(email)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve verification code: %v", err)
	}

	// Compare the saved code with the entered code
	if savedCode == code {
		// Code matches, delete the code from Redis as it is validated
		err = s.redisClient.Delete(email)
		if err != nil {
			return false, fmt.Errorf("failed to delete verified code: %v", err)
		}
		return true, nil
	}

	// Code does not match
	return false, nil
}
