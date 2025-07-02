package seeder

import (
	"github.com/sadamoto/kindle-notifier/internal/db"
	"github.com/sadamoto/kindle-notifier/internal/models"
)

// CategorySeeder handles seeding category data
type CategorySeeder struct{}

// NewCategorySeeder creates a new category seeder
func NewCategorySeeder() *CategorySeeder {
	return &CategorySeeder{}
}

// Seed implements the Seeder interface
func (s *CategorySeeder) Seed() error {
	categories := []models.Category{
		{Name: "ビジネス・経済"},
		{Name: "コンピュータ・IT"},
		{Name: "文学・評論"},
		{Name: "人文・思想"},
		{Name: "社会・政治"},
		{Name: "歴史・地理"},
		{Name: "科学・テクノロジー"},
		{Name: "アート・建築・デザイン"},
		{Name: "趣味・実用"},
		{Name: "教育・自己啓発"},
	}

	for _, category := range categories {
		if err := db.DB.FirstOrCreate(&category, models.Category{Name: category.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}

// Clear implements the Seeder interface
func (s *CategorySeeder) Clear() error {
	return db.DB.Exec("TRUNCATE categories CASCADE").Error
} 