package models

import (
	"time"

	"gorm.io/gorm"
)

// User struct holds data of users
type User struct {
	gorm.Model
	ID             uint       `gorm:"primaryKey"`
	Name           string     `json:"name" binding:"required"`
	Email          string     `json:"email" gorm:"unique" binding:"required"`
	HashedPassword []byte     `json:"hashed_password" binding:"required"`
	UpdatedAt      time.Time  `json:"updated_at"`
	IsOwner        bool       `json:"is_owner"`
	Projects       []*Project `gorm:"many2many:user_projects;"`
}
