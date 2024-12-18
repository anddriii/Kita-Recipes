package domain

import (
	"time"
)

type RecipeAuthor struct {
	ID    int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type RecipeAuthorDetail struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string `json:"name"`
	Photo   string `json:"photo"`
	Recipes Recipe `gorm:"foreignKey:recipe_author_id;references:id" json:"recipes"`
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (u *RecipeAuthorDetail) TableName() string {
	return "recipe_authors"
}
