package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entry struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	Entry_number string    `gorm:"not null;column:entry_number"`
	Content      string    `gorm:"not null;column:content"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at"`
	SearchVector string    `gorm:"type:tsvector"`
}

func (Entry) TableName() string {
	return "entry"
}

func (entry *Entry) BeforeSave(tx *gorm.DB) (err error) {
	entry.SearchVector = formatSearchVector(entry.Entry_number, entry.Content)
	return nil
}

func formatSearchVector(entryNumber, content string) string {
	return fmt.Sprintf("%s %s", entryNumber, content) // Adjust as needed
}

func PopulateSearchVectors(db *gorm.DB) error {
	var entries []Entry
	if err := db.Find(&entries).Error; err != nil {
		return err
	}

	for _, entry := range entries {
		entry.SearchVector = formatSearchVector(entry.Entry_number, entry.Content)
		if err := db.Save(&entry).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *Entry) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return
}

func SearchEntries(db *gorm.DB, pattern string, order string) ([]Entry, error) {
	var entries []Entry

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	searchQuery := fmt.Sprintf("plainto_tsquery('%s')", pattern)

	if err := db.Where("search_vector @@ " + searchQuery).
		Order("string_to_array(entry_number, '.')::int[] " + order).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (e *Entry) FindByID(db *gorm.DB, id uuid.UUID) (*Entry, error) {
	var entry Entry
	if err := db.First(&entry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (e *Entry) Create(db *gorm.DB) error {
	var existing Entry
	if err := db.Where("entry_number = ?", e.Entry_number).First(&existing).Error; err == nil {
		return fmt.Errorf("entry_number %s already exists", e.Entry_number)
	}
	return db.Create(e).Error
}

func (e *Entry) Update(db *gorm.DB) (*Entry, error) {
	if err := db.Omit("ID", "CreatedAt", "Entry_number").Save(e).Error; err != nil {
		return nil, err // Return nil for entry and the error
	}
	return e, nil
}

func (e *Entry) Delete(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&Entry{}, id).Error
}

func FindByID(db *gorm.DB, id uuid.UUID) (*Entry, error) {
	var entry Entry
	if err := db.First(&entry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func FindAll(db *gorm.DB, order string, orderBy string) ([]Entry, error) {
	var entries []Entry
	validOrderBy := map[string]bool{
		"entry_number": true,
		"created_at":   true,
		"updated_at":   true,
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	if !validOrderBy[orderBy] {
		orderBy = "entry_number"
	}

	query := `
		SELECT * 
		FROM entry 
		ORDER BY ` + orderBy

	if orderBy == "entry_number" {
		query += `, string_to_array(entry_number, '.')::int[] ` + order
	} else {
		query += ` ` + order
	}

	if err := db.Raw(query).Scan(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func FindByEntryNumber(db *gorm.DB, pattern string, order string) ([]Entry, error) {
	var entries []Entry
	// entry := "%" + pattern + "%"

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	entry := pattern + "%"

	if err := db.Where("entry_number LIKE ?", entry).
		Order("string_to_array(entry_number, '.')::int[] " + order).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}
