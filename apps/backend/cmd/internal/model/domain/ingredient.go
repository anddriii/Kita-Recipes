package domain

import (
	"time"
)

type Ingredient struct {
	ID    int64  `gorm:"primaryKey:autoIncrement" json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
