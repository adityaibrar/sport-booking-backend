package models

import (
	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID      uint    `gorm:"not null" json:"user_id"`
	VenueID     uint    `gorm:"not null" json:"venue_id"`
	StartTime   string  `gorm:"not null" json:"start_time"`
	Duration    int     `gorm:"not null" json:"duration"` // in hours
	TotalPrice  float64 `gorm:"not null" json:"total_price"`
	Status      string  `gorm:"not null" json:"status"` // pending, confirmed, completed, cancelled
	PaymentQRIS string  `json:"payment_qris"`

	// Add associations
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Venue Venue `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
}

type BookingRequest struct {
	VenueID   uint   `json:"venue_id"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
}

type BookingSearch struct {
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
	Category  string `json:"category"`
}

type BookingResponse struct {
	ID          uint    `json:"id"`
	StartTime   string  `json:"start_time"`
	Duration    int     `json:"duration"`
	TotalPrice  float64 `json:"total_price"`
	Status      string  `json:"status"`
	PaymentQRIS string  `json:"payment_qris,omitempty"`
	CreatedAt   string  `json:"created_at"`
	User        struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"user"`
	Venue struct {
		ID           uint    `json:"id"`
		Name         string  `json:"name"`
		Category     string  `json:"category"`
		PricePerHour float64 `json:"price_per_hour"`
		Description  string  `json:"description"`
	} `json:"venue"`
}
