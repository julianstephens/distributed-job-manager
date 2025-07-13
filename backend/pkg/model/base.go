package model

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string `gorm:"primaryKey;" json:"id,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updatedAt,omitempty"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	guid, err := gonanoid.New()
	if err != nil {
		return err
	}
	base.ID = guid
	return nil
}

