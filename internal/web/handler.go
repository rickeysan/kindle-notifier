package web

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/sadamoto/kindle-notifier/internal/db"
	"github.com/sadamoto/kindle-notifier/internal/models"
)

// Handler handles HTTP requests
type Handler struct {
	templates *template.Template
}

// NewHandler creates a new Handler
func NewHandler() (*Handler, error) {
	// Create template functions
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006/01/02 15:04")
		},
	}

	// Parse templates
	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		templates: tmpl,
	}, nil
}

// HandleIndex handles the index page
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// Get database connection
	database := db.GetDB()
	dbConnected := db.IsConnected()

	// Prepare data structure
	data := struct {
		Books []models.Book
		DBConnected bool
	}{
		Books: []models.Book{}, // Initialize with empty slice
		DBConnected: dbConnected,
	}

	// Only try to fetch books if database is connected
	if dbConnected && database != nil {
		if err := database.Preload("Categories").Order("added_at DESC").Find(&data.Books).Error; err != nil {
			log.Printf("Error fetching books: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Render template
	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
} 