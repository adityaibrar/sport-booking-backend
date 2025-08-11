package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID  uint   `gorm:"not null;index" json:"user_id" validate:"required"`
	VenueID uint   `gorm:"not null;index" json:"venue_id" validate:"required"`
	Rating  int    `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating" validate:"required,min=1,max=5"`
	Comment string `gorm:"size:255" json:"comment" validate:"required"`

	// associations
	User  User  `gorm:"foreignKey:UserID" json:"user"`
	Venue Venue `gorm:"foreignKey:VenueID" json:"venue"`
}

type ReviewRequest struct {
	UserID  uint   `json:"user_id" validate:"required"`
	VenueID uint   `json:"venue_id" validate:"required"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type ReviewResponse struct {
	ID      uint   `json:"id"`
	User    User   `json:"user"`
	Venue   Venue  `json:"venue"`
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewSearchRequest struct {
	VenueCategory VenueCategory `query:"venue_category"`
	Page          int           `query:"page" validate:"min=1"`
	Limit         int           `query:"limit" validate:"min=1,max=100"`
	SortBy        string        `query:"sort_by" validate:"oneof=name price_per_hour created_at"`
	SortOrder     string        `query:"sort_order" validate:"oneof=asc desc"`
}

func (R *Review) ToResponse() ReviewResponse {
	return ReviewResponse{
		ID:      R.ID,
		User:    R.User,
		Venue:   R.Venue,
		Rating:  R.Rating,
		Comment: R.Comment,
	}
}
