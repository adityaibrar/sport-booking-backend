package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the possible roles for users
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// User represents a user in the sports booking system
type User struct {
	gorm.Model
	Name     string   `gorm:"not null;size:100" json:"name" validate:"required,min=2,max=100"`
	Email    string   `gorm:"unique;not null;size:100;index" json:"email" validate:"required,email"`
	Password string   `gorm:"not null;size:255" json:"-" validate:"required,min=8"`
	Phone    string   `gorm:"size:20;index" json:"phone" validate:"omitempty,e164"`
	Address  string   `gorm:"size:500" json:"address" validate:"max=500"`
	Role     UserRole `gorm:"not null;default:'user';index" json:"role" validate:"required,oneof=user admin"`
	IsActive bool     `gorm:"not null;default:true;index" json:"is_active"`

	// Associations
	Bookings []Booking `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bookings,omitempty"`
	Reviews  []Review  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"reviews,omitempty"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone" validate:"omitempty,e164"`
	Address  string `json:"address" validate:"max=500"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserUpdateRequest represents the request payload for updating user profile
type UserUpdateRequest struct {
	Name    string `json:"name" validate:"omitempty,min=2,max=100"`
	Phone   string `json:"phone" validate:"omitempty,e164"`
	Address string `json:"address" validate:"max=500"`
}

// UserResponse represents the response structure for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProfileResponse represents detailed user profile response
type UserProfileResponse struct {
	UserResponse
	TotalBookings     int64   `json:"total_bookings"`
	CompletedBookings int64   `json:"completed_bookings"`
	TotalSpent        float64 `json:"total_spent"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Address:   u.Address,
		Role:      string(u.Role),
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// CanManageBookings checks if user can manage bookings
func (u *User) CanManageBookings() bool {
	return u.IsAdmin()
}
