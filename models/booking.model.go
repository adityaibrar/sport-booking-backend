package models

import (
	"time"

	"gorm.io/gorm"
)

// BookingStatus represents the possible states of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCompleted BookingStatus = "completed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking represents a sports venue booking in the system
type Booking struct {
	gorm.Model
	UserID     uint          `gorm:"not null;index" json:"user_id" validate:"required"`
	VenueID    uint          `gorm:"not null;index" json:"venue_id" validate:"required"`
	StartTime  time.Time     `gorm:"not null;index" json:"start_time" validate:"required"`
	Duration   int           `gorm:"not null" json:"duration" validate:"required,min=1,max=24"` // Duration in hours (1-24)
	TotalPrice float64       `gorm:"not null" json:"total_price" validate:"required,min=0"`
	Status     BookingStatus `gorm:"not null;default:'pending';index" json:"status" validate:"required,oneof=pending confirmed completed cancelled"`
	PaymentUrl string        `gorm:"size:500" json:"payment_url,omitempty"`
	PaymentExpiry time.Time `gorm:"not null" json:"payment_expiry"` // Expiry time for payment link
	Notes      string        `gorm:"size:1000" json:"notes,omitempty"` // Additional booking notes

	// Associations
	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	Venue Venue `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"venue,omitempty"`
}

// BookingRequest represents the request payload for creating a new booking
type BookingRequest struct {
	VenueID   uint      `json:"venue_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	Duration  int       `json:"duration" validate:"required,min=1,max=24"`
	Notes     string    `json:"notes,omitempty" validate:"max=1000"`
}

// BookingUpdateRequest represents the request payload for updating booking status
type BookingUpdateRequest struct {
	Status BookingStatus `json:"status" validate:"required,oneof=pending confirmed completed cancelled"`
	Notes  string        `json:"notes,omitempty" validate:"max=1000"`
}

// BookingSearchRequest represents search parameters for booking queries
type BookingSearchRequest struct {
	UserName  string        `query:"user_name"`
	VenueName string        `query:"venue_name"`
	Status    BookingStatus `query:"status"`
	StartDate string        `query:"start_date"` // Format: YYYY-MM-DD
	EndDate   string        `query:"end_date"`   // Format: YYYY-MM-DD
	MinPrice  float64       `query:"min_price"`
	MaxPrice  float64       `query:"max_price"`
	Page      int           `query:"page" validate:"min=1"`
	Limit     int           `query:"limit" validate:"min=1,max=100"`
	SortBy    string        `query:"sort_by" validate:"oneof=created_at start_time total_price"`
	SortOrder string        `query:"sort_order" validate:"oneof=asc desc"`
}

// BookingAvailabilityRequest represents request for checking venue availability
type BookingAvailabilityRequest struct {
	VenueID   uint   `json:"venue_id" validate:"required"`
	Date      string `json:"date" validate:"required"`       // Format: YYYY-MM-DD
	StartTime string `json:"start_time" validate:"required"` // Format: HH:MM
	Duration  int    `json:"duration" validate:"required,min=1,max=24"`
}

// BookingResponse represents the response structure for booking data
type BookingResponse struct {
	ID                  uint          `json:"id"`
	StartTime           time.Time     `json:"start_time"`
	EndTime             time.Time     `json:"end_time"`
	Duration            int           `json:"duration"`
	TotalPrice          float64       `json:"total_price"`
	Status              string        `json:"status"`
	PaymentUrl          string        `json:"payment_url,omitempty"`
	MidtransToken       string        `json:"midtrans_token,omitempty"`        // Token from Midtrans
	MidtransRedirectURL string        `json:"midtrans_redirect_url,omitempty"` // Redirect URL from Midtrans
	Notes               string        `json:"notes,omitempty"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	User                UserResponse  `json:"user"`
	Venue               VenueResponse `json:"venue"`
}

// BookingListResponse represents paginated booking list response
type BookingListResponse struct {
	Data       []BookingResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// BookingStatsResponse represents booking statistics
type BookingStatsResponse struct {
	TotalBookings     int64   `json:"total_bookings"`
	PendingBookings   int64   `json:"pending_bookings"`
	ConfirmedBookings int64   `json:"confirmed_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	MonthlyRevenue    float64 `json:"monthly_revenue"`
}

// GetEndTime calculates the end time based on start time and duration
func (b *Booking) GetEndTime() time.Time {
	return b.StartTime.Add(time.Duration(b.Duration) * time.Hour)
}

// IsActive checks if the booking is currently active
func (b *Booking) IsActive() bool {
	now := time.Now()
	return b.StartTime.Before(now) && b.GetEndTime().After(now) && b.Status == BookingStatusConfirmed
}

// CanBeCancelled checks if the booking can be cancelled
func (b *Booking) CanBeCancelled() bool {
	return b.Status == BookingStatusPending || (b.Status == BookingStatusConfirmed && time.Until(b.StartTime) > time.Hour)
}

// ToResponse converts Booking model to BookingResponse
func (b *Booking) ToResponse(midtransToken string, midtransRedirectURL string) BookingResponse {
	response := BookingResponse{
		ID:         b.ID,
		StartTime:  b.StartTime,
		EndTime:    b.GetEndTime(),
		Duration:   b.Duration,
		TotalPrice: b.TotalPrice,
		Status:     string(b.Status),
		PaymentUrl: b.PaymentUrl,
		Notes:      b.Notes,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
	}

	// Add Midtrans token and redirect URL if provided
	if midtransToken != "" {
		response.MidtransToken = midtransToken
	}
	if midtransRedirectURL != "" {
		response.MidtransRedirectURL = midtransRedirectURL
	}

	// Convert user data if loaded
	if b.User.ID != 0 {
		response.User = b.User.ToResponse()
	}

	// Convert venue data if loaded
	if b.Venue.ID != 0 {
		response.Venue = b.Venue.ToResponse()
	}

	return response
}
