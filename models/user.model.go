package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Phone    int32  `json:"phone"`
	Address  string `json:"address"`
	Role     string `gorm:"default:user" json:"role"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    int32 `json:"phone"`
	Address  string `json:"address"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}
