package config

import (
	"log"
	"sport-booking-backend/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		log.Printf("Failed to auto migrate :%v", err)
		return err
	}
	log.Println("Database migratrion completed succesfully")
	return nil
}
