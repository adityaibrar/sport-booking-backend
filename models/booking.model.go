package models

import (
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID      float64   `gorm:"not null" json:"user_id"`
	VenueID     uint      `gorm:"not null" json:"venue_id"`
	StartTime   string `gorm:"not null" json:"start_time"`
	Duration    int       `gorm:"not null" json:"duration"` // in hours
	TotalPrice  float64   `gorm:"not null" json:"total_price"`
	Status      string    `gorm:"not null" json:"status"` // pending, confirmed, completed, cancelled
	PaymentQRIS string    `json:"payment_qris"`
}

type BookingRequest struct {
	VenueID   uint   `json:"venue_id"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
}
