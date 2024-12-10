package middleware

import (
	"errors"
	"online-questionnaire/internal/logger"
	"online-questionnaire/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.GetLogger().Error("Could not get authorization token", errors.New("Authorization header is missing"), logger.Logctx{}, "")

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is missing"})
		}

		// Bearer token validation
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			logger.GetLogger().Error("Invalit format.", errors.New("Invalid Authorization header format"), logger.Logctx{Data: map[string]interface{}{
				"App":    c.App,
				"route":  c.Route,
				"Method": c.Method,
			}}, "")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
		}

		// Parse the token
		tokenString := tokenParts[1]
		token, err := jwt.ParseWithClaims(tokenString, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.GetLogger().Error("Unsupported Method", errors.New("unexpected signing method"), logger.Logctx{Data: map[string]interface{}{
					"App":    c.App,
					"route":  c.Route,
					"Method": c.Method,
				}}, "")
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logger.GetLogger().Error("Invalid token", errors.New("Token is invalid or expiered"), logger.Logctx{Data: map[string]interface{}{
				"App":    c.App,
				"route":  c.Route,
				"Method": c.Method,
			}}, "")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		// Extract claims
		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok {
			logger.GetLogger().Error("Invalid token claims", errors.New("Invalid token claims"), logger.Logctx{Data: map[string]interface{}{
				"App":    c.App,
				"route":  c.Route,
				"Method": c.Method,
			}}, "")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// Attach claims to context
		c.Locals("user_id", claims.Username)
		c.Locals("role", claims.Role)

		// Proceed to the next middleware or handler
		return c.Next()
	}
}
