package seeder

import (
	"time"

	"github.com/sadamoto/kindle-notifier/internal/db"
	"github.com/sadamoto/kindle-notifier/internal/models"
)

// BookSeeder handles seeding book data
type BookSeeder struct{}

// NewBookSeeder creates a new book seeder
func NewBookSeeder() *BookSeeder {
	return &BookSeeder{}
}

// Seed implements the Seeder interface
func (s *BookSeeder) Seed() error {
	books := []models.Book{
		{
			ASIN:       "B0XXXXX1",
			Title:      "Go言語プログラミング入門",
			Author:     "山田太郎",
			ImageURL:   "https://example.com/image1.jpg",
			ProductURL: "https://amazon.co.jp/dp/B0XXXXX1",
			IsKU:       true,
			AddedAt:    time.Now(),
			Description: "Go言語の基礎から応用までを解説する入門書です。",
		},
		{
			ASIN:       "B0XXXXX2",
			Title:      "実践マイクロサービス",
			Author:     "鈴木一郎",
			ImageURL:   "https://example.com/image2.jpg",
			ProductURL: "https://amazon.co.jp/dp/B0XXXXX2",
			IsKU:       true,
			AddedAt:    time.Now(),
			Description: "マイクロサービスアーキテクチャの実践的な解説書です。",
		},
	}

	// Create books
	for _, book := range books {
		if err := db.DB.FirstOrCreate(&book, models.Book{ASIN: book.ASIN}).Error; err != nil {
			return err
		}

		// Find corresponding categories
		var categories []models.Category
		if book.Title == "Go言語プログラミング入門" {
			if err := db.DB.Where("name IN ?", []string{"コンピュータ・IT", "教育・自己啓発"}).Find(&categories).Error; err != nil {
				return err
			}
		} else {
			if err := db.DB.Where("name IN ?", []string{"コンピュータ・IT", "ビジネス・経済"}).Find(&categories).Error; err != nil {
				return err
			}
		}

		// Associate categories with the book
		if err := db.DB.Model(&book).Association("Categories").Replace(categories); err != nil {
			return err
		}
	}

	return nil
}

// Clear implements the Seeder interface
func (s *BookSeeder) Clear() error {
	return db.DB.Exec("TRUNCATE books CASCADE").Error
} 