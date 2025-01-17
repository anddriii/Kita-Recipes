package domain

import "time"

type Login struct {
	ID        int    `gorm:"primaryKey:autoIncrement" json:"id"`
	AuthorId  int    `gorm:"not null;unique"`
	Username  string `gorm:"not null;unique"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    RecipeAuthor `gorm:"foreignKey:AuthorId"`
}
