package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project model, containes all project fields.
type Project struct {
	gorm.Model
	ID              uuid.UUID         `gorm:"primary_key; unique; type:uuid; column:id"`
	Name            string            `json:"name" binding:"required"`
	EnvironmentName string            `json:"environment_name"`
	Team            []*User           `gorm:"many2many:project_team;default:nil"`
	Owner           uuid.UUID         // Foreign key referencing User's ID field
	Keys            []*EnvironmentKey `gorm:"default:nil"`
}

// Env keys model, containes all project keys.
type EnvironmentKey struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id"`
	ProjectID uuid.UUID
	Key       string `gorm:"unique;"`
	Value     []byte
}

// BeforeCreate generates a new uuid.
func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	p.ID = id
	return
}

// BeforeCreate generates a new uuid.
func (env *EnvironmentKey) BeforeCreate(*gorm.DB) (err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error, ", err)
		return err
	}

	env.ID = id
	return
}
