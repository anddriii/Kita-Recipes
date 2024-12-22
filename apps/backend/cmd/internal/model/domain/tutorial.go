package domain

import (
	"time"
)

type RecipeTutorial struct {
	ID       int64 `gorm:"primaryKey:autoIncrement"`
	Name     string
	RecipeId int64
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
