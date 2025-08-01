package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp_id":  time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString, err
}

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No Token Provided",
		})
	}

	var tokenString string
	fmt.Sscanf(authHeader, "Beare %s", &tokenString)

	if tokenString == "" {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": "Token format is invalid",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is invalid",
		})
	}
	return c.Next()
}
