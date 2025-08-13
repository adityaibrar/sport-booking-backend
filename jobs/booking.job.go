package jobs

import (
	"log"
	"sport-booking-backend/models"
	"time"

	"gorm.io/gorm"
)

func CancelExpiredBookings(db *gorm.DB) {
	now := time.Now()
	var expiredBookings []models.Booking
	if err := db.Where("status = ? AND payment_expiry <= ?", models.BookingStatusPending, now).Find(&expiredBookings).Error; err != nil {
		log.Printf("Failed to fetch expired bookings: %v", err)
		return
	}
	for _, booking := range expiredBookings {
		booking.Status = models.BookingStatusCancelled
		if err := db.Save(&booking).Error; err != nil {
			log.Printf("Failed to cancel booking %d: %v", booking.ID, err)
			return
		}
	}
}
