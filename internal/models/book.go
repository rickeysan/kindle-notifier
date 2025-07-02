package models

import (
	"time"
)

// Book represents a Kindle Unlimited book
type Book struct {
	ID          uint      `gorm:"primaryKey"`
	ASIN        string    `gorm:"uniqueIndex;not null"`
	Title       string    `gorm:"not null"`
	Author      string    `gorm:"not null"`
	ImageURL    string    `gorm:"not null"`
	ProductURL  string    `gorm:"not null"`
	IsKU        bool      `gorm:"not null"`
	AddedAt     time.Time `gorm:"not null"`
	Description string
	Categories  []Category `gorm:"many2many:book_categories;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Category represents a book category
type Category struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"uniqueIndex;not null"`
	Books []Book `gorm:"many2many:book_categories;"`
} 