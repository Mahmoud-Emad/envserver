package models

import (
	"gorm.io/gorm"
)

// Project model, containes all project fields.
type Project struct {
	gorm.Model
	ID              uint              `gorm:"primaryKey"`
	Name            string            `json:"name" binding:"required"`
	EnvironmentName string            `json:"environment_name"` // e.g. test, dev, production.
	Team            []*User           `gorm:"many2many:project_team;default:nil"`
	Owner           uint              // Foreign key referencing User's ID field
	Keys            []*EnvironmentKey `gorm:"default:nil"`
}

// Env keys model, containes all project keys.
type EnvironmentKey struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	ProjectID uint
	Key       string
	Value     []byte
}
