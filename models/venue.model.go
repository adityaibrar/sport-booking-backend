package models

import (
	"time"

	"gorm.io/gorm"
)

// VenueCategory represents the categories of sports venues
type VenueCategory string

const (
	VenueCategoryFootball   VenueCategory = "football"
	VenueCategoryBasketball VenueCategory = "basketball"
	VenueCategoryBadminton  VenueCategory = "badminton"
	VenueCategoryTennis     VenueCategory = "tennis"
	VenueCategoryVolleyball VenueCategory = "volleyball"
	VenueCategorySwimming   VenueCategory = "swimming"
	VenueCategoryFutsal     VenueCategory = "futsal"
	VenueCategoryOther      VenueCategory = "other"
)

// Venue represents a sports venue in the system
type Venue struct {
	gorm.Model
	Name         string        `gorm:"not null;size:100;index" json:"name" validate:"required,min=2,max=100"`
	Category     VenueCategory `gorm:"not null;size:50;index" json:"category" validate:"required,oneof=football basketball badminton tennis volleyball swimming futsal other"`
	PricePerHour float64       `gorm:"not null" json:"price_per_hour" validate:"required,min=0"`
	Description  string        `gorm:"not null;size:1000" json:"description" validate:"required,max=1000"`
	Location     string        `gorm:"size:500" json:"location" validate:"max=500"`
	Facilities   string        `gorm:"size:1000" json:"facilities" validate:"max=1000"` // JSON string of facilities
	Images       string        `gorm:"size:2000" json:"images" validate:"max=2000"`     // JSON string of image URLs
	IsActive     bool          `gorm:"not null;default:true;index" json:"is_active"`
	OpenTime     string        `gorm:"size:5" json:"open_time" validate:"required"`  // Format: HH:MM (24h)
	CloseTime    string        `gorm:"size:5" json:"close_time" validate:"required"` // Format: HH:MM (24h)
	Capacity     int           `gorm:"not null;default:1" json:"capacity" validate:"min=1"`

	// Associations
	Bookings []Booking `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bookings,omitempty"`
}

// VenueRequest represents the request payload for creating/updating a venue
type VenueRequest struct {
	Name         string        `json:"name" validate:"required,min=2,max=100"`
	Category     VenueCategory `json:"category" validate:"required,oneof=football basketball badminton tennis volleyball swimming futsal other"`
	PricePerHour float64       `json:"price_per_hour" validate:"required,min=0"`
	Description  string        `json:"description" validate:"required,max=1000"`
	Location     string        `json:"location" validate:"max=500"`
	Facilities   []string      `json:"facilities" validate:"max=20,dive,max=100"`
	Images       []string      `json:"images" validate:"max=10,dive,url"`
	OpenTime     string        `json:"open_time" validate:"required,len=5"`  // Format: HH:MM
	CloseTime    string        `json:"close_time" validate:"required,len=5"` // Format: HH:MM
	Capacity     int           `json:"capacity" validate:"min=1"`
}

// VenueSearchRequest represents search parameters for venue queries
type VenueSearchRequest struct {
	Name      string        `query:"name"`
	Category  VenueCategory `query:"category"`
	MinPrice  float64       `query:"min_price"`
	MaxPrice  float64       `query:"max_price"`
	Location  string        `query:"location"`
	IsActive  *bool         `query:"is_active"`
	Page      int           `query:"page" validate:"min=1"`
	Limit     int           `query:"limit" validate:"min=1,max=100"`
	SortBy    string        `query:"sort_by" validate:"oneof=name price_per_hour created_at"`
	SortOrder string        `query:"sort_order" validate:"oneof=asc desc"`
}

// VenueResponse represents the response structure for venue data
type VenueResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	PricePerHour float64   `json:"price_per_hour"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	Facilities   []string  `json:"facilities"`
	Images       []string  `json:"images"`
	IsActive     bool      `json:"is_active"`
	OpenTime     string    `json:"open_time"`
	CloseTime    string    `json:"close_time"`
	Capacity     int       `json:"capacity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// VenueDetailResponse represents detailed venue response with statistics
type VenueDetailResponse struct {
	VenueResponse
	TotalBookings      int64    `json:"total_bookings"`
	ActiveBookings     int64    `json:"active_bookings"`
	TotalRevenue       float64  `json:"total_revenue"`
	AverageRating      float64  `json:"average_rating"`
	BookingRate        float64  `json:"booking_rate"` // Percentage of time booked
	AvailableTimeSlots []string `json:"available_time_slots,omitempty"`
}

// VenueListResponse represents paginated venue list response
type VenueListResponse struct {
	Data       []VenueResponse `json:"data"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// VenueAvailabilityResponse represents venue availability information
type VenueAvailabilityResponse struct {
	VenueID          uint     `json:"venue_id"`
	Date             string   `json:"date"`
	IsAvailable      bool     `json:"is_available"`
	AvailableSlots   []string `json:"available_slots"`
	UnavailableSlots []string `json:"unavailable_slots"`
}

// ToResponse converts Venue model to VenueResponse
func (v *Venue) ToResponse() VenueResponse {
	return VenueResponse{
		ID:           v.ID,
		Name:         v.Name,
		Category:     string(v.Category),
		PricePerHour: v.PricePerHour,
		Description:  v.Description,
		Location:     v.Location,
		Facilities:   []string{}, // Parse from JSON string if needed
		Images:       []string{}, // Parse from JSON string if needed
		IsActive:     v.IsActive,
		OpenTime:     v.OpenTime,
		CloseTime:    v.CloseTime,
		Capacity:     v.Capacity,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
	}
}

// IsOpenAt checks if venue is open at the given time
func (v *Venue) IsOpenAt(timeStr string) bool {
	return timeStr >= v.OpenTime && timeStr <= v.CloseTime
}

// GetOperatingHours returns the operating hours as a slice of hour strings
func (v *Venue) GetOperatingHours() []string {
	var hours []string
	// Implementation would parse open/close times and generate hourly slots
	// This is a simplified version
	for i := 0; i < 24; i++ {
		hourStr := time.Date(0, 0, 0, i, 0, 0, 0, time.UTC).Format("15:04")
		if v.IsOpenAt(hourStr) {
			hours = append(hours, hourStr)
		}
	}
	return hours
}
