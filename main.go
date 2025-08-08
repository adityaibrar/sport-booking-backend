package main

import (
	"log"
	"os"
	"sport-booking-backend/config"
	"sport-booking-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	dsn := "root:abdillah24@tcp(127.0.0.1:3306)/sport_booking?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	if err := config.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	routes.SetupRoutes(app, db)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
