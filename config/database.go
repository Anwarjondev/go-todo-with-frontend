package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "go-todo")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Retry database connection
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		// Open database connection
		var err error
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: failed to connect to database: %v. Retrying...", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Test the connection
		if err = DB.Ping(); err != nil {
			log.Printf("Attempt %d: failed to ping database: %v. Retrying...", i+1, err)
			DB.Close() // Close the potentially invalid connection
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Successfully connected to database")
		// If connection is successful, break the loop
		break
	}

	// Check if the connection was successful after retries
	if DB == nil || DB.Ping() != nil {
		return fmt.Errorf("failed to connect to database after %d retries", maxRetries)
	}

	// Create todos table if it doesn't exist
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			item TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			user_id INTEGER NOT NULL
		);
	`

	_, err := DB.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create todos table: %v", err)
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// Helper function to get environment variables with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
