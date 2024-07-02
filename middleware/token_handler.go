package middleware

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/kuusminustwo/go_rest_api/database"
	"github.com/kuusminustwo/go_rest_api/model"
)

var jwtSecret = []byte("your-secret-key")

// VerifyToken verifies the JWT token and authorizes access to protected endpoints
func TokenAuthMiddleware(c *fiber.Ctx) error {
	// fmt.Println("TokenAuthMiddleware called")
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	tokenString := extractTokenFromHeader(authHeader)

	fmt.Println("Token extracted: " + tokenString)
	// Parse JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    err.Error(),
		})
	}

	// Verify token validity
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
			"data":    nil,
		})
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token claims",
			"data":    nil,
		})
	}

	// Optionally, fetch user details from database using claims
	userID := claims["userID"].(string) // Example of extracting userID claim

	// Fetch user details from database using userID
	db := database.DB.Db
	var user model.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    nil,
		})
	}

	// Store user information in context for further access if needed
	c.Locals("user", user)

	// Proceed to the next middleware or handler
	return c.Next()
}

// extractTokenFromHeader extracts the token from the Authorization header
func extractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	if len(tokenParts) != 2 {
		return ""
	}
	return tokenParts[1]
}
