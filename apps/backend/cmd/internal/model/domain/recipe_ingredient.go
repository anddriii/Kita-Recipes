package domain

import (
	"time"
)

type RecipeIngredient struct {
	ID           int64 `gorm:"primaryKey:autoIncrement" json:"id"`
	IngredientId int64 `gorm:"foreignKey:ingredient_id;references:id"`
	RecipeId     int64 `gorm:"foreignKey:recipe_id;references:id"`
	// DeletedAt    gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
