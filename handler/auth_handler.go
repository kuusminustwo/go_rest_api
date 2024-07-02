package handler

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kuusminustwo/go_rest_api/database"
	"github.com/kuusminustwo/go_rest_api/model"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure secret key

// LoginInput represents the structure for incoming login requests
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the structure for the login response
type LoginResponse struct {
	Token string `json:"token"`
}

// generateJWTToken generates a new JWT token for the given user ID
func generateJWTToken(userID uuid.UUID) (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Create claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID.String()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiry in 24 hours

	// Generate encoded token and return it
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Login handles user authentication and returns a JWT token upon successful login
func Login(c *fiber.Ctx) error {
	db := database.DB.Db

	// Parse request body into LoginInput struct
	var loginData LoginInput
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	// Fetch user from database
	var user model.User
	if err := db.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
			"data":    nil,
		})
	}

	// Compare hashed password with the password provided
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credentials",
			"data":    nil,
		})
	}

	// Generate JWT token
	token, err := generateJWTToken(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to generate token",
			"data":    nil,
		})
	}

	// Return token as response
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data":    LoginResponse{Token: token},
	})
}
