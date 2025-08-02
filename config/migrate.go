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
		log.Printf("Failed to auto migrate :%v", err)
		return err
	}

	// err = db.Exec(`
	// 	ALTER TABLE bookings
	// 	ADD CONSTRAINT fk_bookings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	// 	ADD CONSTRAINT fk_bookings_venue FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE CASCADE
	// `).Error
	// if err != nil {
	// 	log.Printf("Failed to add foreign key constraints: %v", err)
	// 	return err
	// }

	log.Println("Database migratrion completed succesfully")
	return nil
}
