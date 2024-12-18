package domain

import (
	"time"
)

type Categories struct {
	ID   int64  `gorm:"primaryKey:autoIncrement" json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Icon string `json:"icon"`
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type CategoryDetail struct {
	ID   int64  `gorm:"primaryKey:autoIncrement" json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Icon string `json:"icon"`
	// DeletedAt gorm.DeletedAt `gorm:"autoCreateTime"` //soft deleted otomatis dari GORM
	UpdatedAt time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	//relasi hasMany
	Recipes []Recipe `gorm:"foreignKey:category_id;references:id"`
}

func (u *CategoryDetail) TableName() string {
	return "categories"
}
