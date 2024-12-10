package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"online-questionnaire/internal/utils"
	"strings"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
		}

		// Bearer token validation
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
		}

		// Parse the token
		tokenString := tokenParts[1]
		token, err := jwt.ParseWithClaims(tokenString, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Extract claims
		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Attach claims to context
		c.Locals("user_id", claims.Username)
		c.Locals("role", claims.Role)

		// Proceed to the next middleware or handler
		return c.Next()
	}
}
