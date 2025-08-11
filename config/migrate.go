package config

import (
	"log"
	"sport-booking-backend/models"
	"sport-booking-backend/seeders"

	"gorm.io/gorm"
)

func ResetAndMigrate(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
		&models.Venue{},
		&models.Booking{},
		&models.Review{},
	}

	if err := db.Migrator().DropTable(models...); err != nil {
		log.Printf("Failed to drop tables: %v", err)
		return err
	}

	if err := db.AutoMigrate(models...); err != nil {
		log.Printf("Failed to auto migrate: %v", err)
		return err
	}

	log.Println("Running seeders...")
	if err := seeders.SeedAdminUser(db); err != nil {
		return err
	}

	log.Println("Database reset and migration completed successfully")
	return nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Venue{},
		&models.Booking{},
		&models.Review{},
	)

	if err != nil {
		log.Printf("Failed to auto migrate: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
