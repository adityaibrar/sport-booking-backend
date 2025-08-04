package config

import (
	"log"
	"sport-booking-backend/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Venue{},
		&models.Booking{},
	)

	if err != nil {
		log.Printf("Failed to auto migrate: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}