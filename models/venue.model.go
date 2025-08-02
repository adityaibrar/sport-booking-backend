package models

import "gorm.io/gorm"

type Venue struct {
	gorm.Model
	Name         string  `gorm:"not null" json:"name"`
	Category     string  `gorm:"not null" json:"category"`
	PricePerHour float64 `gorm:"not null" json:"price_per_hour"`
	Description  string  `gorm:"not null" json:"description"`
}

type VenueRequest struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	PricePerHour float64 `json:"price_per_hour"`
	Description  string  `json:"description"`
}
