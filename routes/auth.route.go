package routes

import (
	"sport-booking-backend/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupAuthRoutes(version fiber.Router, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	version.Post("/register", authController.Register)
	version.Post("/login", authController.Login)
}