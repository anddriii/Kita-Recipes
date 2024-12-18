package domain

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        int64  `gorm:"primaryKey:autoIncrement" json:"id"`
	Photo     string `json:"photo"`
	RecipeId  int64
	DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time      `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}

func (u *Photo) TableName() string {
	return "recipe_photos"
}
