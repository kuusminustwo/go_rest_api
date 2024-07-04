package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kuusminustwo/go_rest_api/database"
	"github.com/kuusminustwo/go_rest_api/model"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)

	// Parse request body into User struct
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    err,
		})
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to hash password",
			"data":    err,
		})
	}
	user.Password = string(hashedPassword)

	// Store user in database
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User has been created",
		"data":    user,
	})
}
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.Db
	var users []model.User
	db.Find(&users)
	if len(users) == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Users not found", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Users Found", "data": users})
}

func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.Db
	id := c.Params("id")
	var user model.User
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

func UpdateUser(c *fiber.Ctx) error {
	type updateUser struct {
		Username string `json:"username"`
	}
	db := database.DB.Db
	var user model.User
	id := c.Params("id")
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	user.Username = updateUserData.Username
	db.Save(&user)
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "users Found", "data": user})
}

func DeleteUserByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var user model.User
	id := c.Params("id")
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
	err := db.Delete(&user, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete user", "data": nil})
	}
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User deleted"})
}

func GetUserProfile(c *fiber.Ctx) error {
	// Access authenticated user from context
	user := c.Locals("user").(model.User)

	// Preload student data
	db := database.DB.Db
	var student model.Student
	if err := db.First(&student, "user_id = ?", user.ID).Error; err != nil {
		// Handle error if student data not found
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch student data",
			"data":    err.Error(),
		})
	}

	// Create response structure including both user and student data
	response := fiber.Map{
		"status":  "success",
		"message": "User profile retrieved",
		"data": fiber.Map{
			"user":    user,
			"student": student,
		},
	}

	// Return the combined data as JSON response
	return c.JSON(response)
}

func GetStudentEnrollments(c *fiber.Ctx) error {
	// Extract the user from the context
	user := c.Locals("user").(model.User)

	// Find the student associated with the user
	db := database.DB.Db
	var student model.Student
	if err := db.Where("user_id = ?", user.ID).First(&student).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	// Get the enrollments for the student
	enrollments, err := model.GetEnrollmentsByStudentID(db, student.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve enrollments",
		})
	}

	return c.JSON(fiber.Map{
		"student":     student,
		"enrollments": enrollments,
	})
}
