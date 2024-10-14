package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUID struct {
	Id string `gorm:"primarykey"`
}

func (u *UUID) BeforeCreate(tx *gorm.DB) error {
	if len(u.Id) == 0 {
		u.Id = uuid.NewString()
	}
	return nil
}
