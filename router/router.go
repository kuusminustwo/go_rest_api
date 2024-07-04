package router

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kuusminustwo/go_rest_api/handler"
	"github.com/kuusminustwo/go_rest_api/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	v1 := api.Group("/user")

	fmt.Println("Setting up routes...")

	v1.Get("/", handler.GetAllUsers)
	v1.Use(middleware.TokenAuthMiddleware) // Middleware to verify token for all routes in this groupasdasd

	fmt.Println("TokenAuthMiddleware applied to /api/user routes")

	v1.Get("/:id", handler.GetSingleUser)
	v1.Post("/", handler.CreateUser)
	v1.Put("/:id", handler.UpdateUser)
	v1.Delete("/:id", handler.DeleteUserByID)
	api.Post("/login", handler.Login)
	api.Get("/profile", middleware.TokenAuthMiddleware, handler.GetUserProfile)
	api.Get("/tod", middleware.TokenAuthMiddleware, handler.GetUserProfile)
	api.Get("/dun", middleware.TokenAuthMiddleware, handler.GetStudentEnrollments)
	api.Get("/hicheel", middleware.TokenAuthMiddleware, handler.GetUserProfile)
}
