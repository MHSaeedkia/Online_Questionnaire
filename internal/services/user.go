package services

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/smtp"
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

// GenerateOTP to generate a 5-digit OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%05d", rand.Intn(100000)) // 5-digit number
	return otp
}

// SendEmailOTP to send the OTP via email
func SendEmailOTP(otp string, userEmail string) error {
	from := "saharhallaji.dev@gmail.com"
	password := "your-email-password" // Use an app-specific password or OAuth for security
	to := []string{userEmail}
	subject := "Your One-Time Password (OTP)"
	body := fmt.Sprintf("Your OTP is: %s", otp)

	msg := []byte("From: " + from + "\r\n" +
		"To: " + userEmail + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, to, msg)

	return err
}

// GenerateCode to generate and store OTP in Redis for 2FA
func GenerateCode() string {
	// Create a new random generator with a unique seed
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Generate a 5-digit random number
	code := r.Intn(100000) // Generates a number between 0 and 99999

	// Ensure it's always 5 digits (padding with leading zeros if necessary)
	return fmt.Sprintf("%05d", code)
}

// ValidateOTP to validate the OTP provided by the user
func ValidateOTP(nationalID string, inputOTP string) (bool, error) {
	// Get the stored OTP from Redis
	storedOTP, err := redisClient.Get(ctx, nationalID).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("OTP expired or not found")
	} else if err != nil {
		return false, err
	}

	// Compare the stored OTP with the input OTP
	if storedOTP == inputOTP {
		// OTP matches, delete it from Redis
		redisClient.Del(ctx, nationalID)
		return true, nil
	}
	return false, fmt.Errorf("incorrect OTP")
}
