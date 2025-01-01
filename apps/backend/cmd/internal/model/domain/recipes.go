package domain

import (
	"time"
)

type Recipe struct {
	ID             int64      `gorm:"primaryKey:autoIncrement" json:"id"`
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	Thumbnail      string     `json:"thumbnail"`
	About          string     `json:"about"`
	CategoryId     int64      `gorm:"column:category_id" json:"category_id"`
	Category       Categories `gorm:"foreignKey:CategoryId;references:id"`
	RecipeAuthorId int64
	RecipeAuthor   RecipeAuthor     `gorm:"foreignKey:recipe_author_id;references:id"`
	RecipePhoto    []Photo          `gorm:"foreignKey:recipe_id;references:id"`
	UrlVideo       string           `json:"url_video"`
	UrlFile        string           `json:"url_file"`
	RecipeTutorial []RecipeTutorial `gorm:"foreignKey:recipe_id;references:id"`
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type RecipeDetail struct {
	ID             int64            `gorm:"primaryKey:autoIncrement" json:"id"`
	Name           string           `json:"name"`
	Slug           string           `json:"slug"`
	Thumbnail      string           `json:"thumbnail"`
	About          string           `json:"about"`
	CategoryId     int64            `form:"category_id" json:"category_id"`
	Category       Categories       `gorm:"foreignKey:CategoryId;references:id"`
	RecipeAuthorId int64            `gorm:"column:recipe_author_id"`
	RecipeAuthor   RecipeAuthor     `gorm:"foreignKey:recipe_author_id;references:id"`
	UrlVideo       string           `json:"url_video"`
	UrlFile        string           `json:"url_file"`
	RecipeTutorial []RecipeTutorial `gorm:"foreignKey:recipe_id;references:id"` //many to many relasi
	Ingredients    []Ingredient     `gorm:"many2many:recipe_ingredients;foreignKey:id;joinForeignKey:recipe_id;References:ID;joinReferences:ingredient_id"`
	RecipePhoto    []Photo          `gorm:"foreignKey:recipe_id;references:id"`
	// DeletedAt      gorm.DeletedAt   `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (u *RecipeDetail) TableName() string {
	return "recipes"
}
