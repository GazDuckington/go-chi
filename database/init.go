package database

import (
	model "bobot/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DatabaseURL string
var DB *gorm.DB
var err error

func autoMigrate() {

	// entry table
	err = DB.AutoMigrate(&model.Entry{}, &model.User{})
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}
	log.Printf("model: %s created", model.Entry{}.TableName())

	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_search_vector ON entry USING GIN(search_vector);").Error; err != nil {
		log.Fatalf("failed to create search index: %v", err)
	}
	log.Print("search vector index created")

	if err := model.PopulateSearchVectors(DB); err != nil {
		log.Fatalf("failed to populate search vectors: %v", err)
	}
	log.Print("search vectors populated")
}

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access environment variables
	pgHost := os.Getenv("PGHOST")
	pgUser := os.Getenv("PGUSER")
	pgPass := os.Getenv("PGPASSWORD")
	pgDb := os.Getenv("PGDATABASE")

	DatabaseURL = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=require", pgUser, pgPass, pgHost, pgDb)
	DB, err = gorm.Open(postgres.Open(DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Printf("connected to database %s", DatabaseURL)

	autoMigrate()

}
