package routes

import (
	"sport-booking-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/register", authController.Register)
	v1.Post("/login", authController.Login)
}
