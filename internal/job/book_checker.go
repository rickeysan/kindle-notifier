package job

import (
	"log"
	"time"

	"github.com/sadamoto/kindle-notifier/internal/api"
	"github.com/sadamoto/kindle-notifier/internal/db"
	"github.com/sadamoto/kindle-notifier/internal/models"
)

// BookChecker handles checking for new Kindle Unlimited books
type BookChecker struct {
	client *api.PAAPIClient
}

// NewBookChecker creates a new BookChecker instance
func NewBookChecker(client *api.PAAPIClient) *BookChecker {
	return &BookChecker{
		client: client,
	}
}

// Run executes the book checking process
func (c *BookChecker) Run() error {
	log.Println("Starting to check for new Kindle Unlimited books...")
	
	// Search for Kindle Unlimited books
	response, err := c.client.SearchKindleUnlimitedBooks()
	if err != nil {
		return err
	}

	// Process each book
	for _, item := range response.ItemsResult.Items {
		// Check if book already exists
		var existingBook models.Book
		result := db.DB.Where("asin = ?", item.ASIN).First(&existingBook)
		if result.Error == nil {
			// Book already exists, skip it
			continue
		}

		// Get author name
		var author string
		if len(item.ItemInfo.ByLineInfo.Contributors) > 0 {
			author = item.ItemInfo.ByLineInfo.Contributors[0].Name
		}

		// Create new book
		newBook := models.Book{
			ASIN:       item.ASIN,
			Title:      item.ItemInfo.Title.DisplayValue,
			Author:     author,
			ImageURL:   item.Images.Primary.Large.URL,
			ProductURL: item.DetailPageURL,
			IsKU:      true,
			AddedAt:    time.Now(),
		}

		// Save to database
		if err := db.DB.Create(&newBook).Error; err != nil {
			log.Printf("Error saving book %s: %v", item.ASIN, err)
			continue
		}

		log.Printf("Added new book: %s by %s", newBook.Title, newBook.Author)
	}

	log.Println("Finished checking for new Kindle Unlimited books")
	return nil
} 