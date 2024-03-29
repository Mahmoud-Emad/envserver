package models

import (
	"time"

	"gorm.io/gorm"
)

// User struct holds data of users
type User struct {
	gorm.Model
	ID             int        `gorm:"primaryKey"`
	FirstName      string     `json:"first_name" binding:"required"`
	LastName       string     `json:"last_name" binding:"required"`
	Email          string     `json:"email" gorm:"unique" binding:"required"`
	HashedPassword []byte     `json:"hashed_password" binding:"required"`
	UpdatedAt      time.Time  `json:"updated_at"`
	IsOwner        bool       `json:"is_owner"`
	Projects       []*Project `gorm:"many2many:user_projects;"`
}
