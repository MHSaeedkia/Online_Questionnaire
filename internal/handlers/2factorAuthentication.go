package handlers

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"math/rand"
	"net/smtp"
	"time"
)

type TwoFactorHandler struct {
	RedisClient *redis.Client
}

func NewTwoFactorHandler(redisClient *redis.Client) *TwoFactorHandler {
	return &TwoFactorHandler{RedisClient: redisClient}
}

func (h *TwoFactorHandler) GenerateCode(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}

	req := new(Request)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Generate a 5-digit code
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%05d", rand.Intn(100000))

	// Save the code in Redis with a 5-minute expiration
	ctx := context.Background()
	err := h.RedisClient.Set(ctx, req.Email, code, 5*time.Minute).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save the code"})
	}

	// Send the code via email
	err = sendEmail(req.Email, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not send email"})
	}

	return c.JSON(fiber.Map{"message": "Code sent successfully"})
}

func sendEmail(to, code string) error {
	from := "saharhallaji.dev@gmail.com"
	password := "your-email-password" // Replace with your email password or app-specific password

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(fmt.Sprintf("Subject: Your 2FA Code\n\nYour verification code is: %s", code))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	return err
}
