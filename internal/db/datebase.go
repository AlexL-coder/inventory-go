package db

import (
	"awesomeProject1/ent"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

var Client *ent.Client

func InitDB(Log *logrus.Logger) {
	Log.Info("Attempting to connect to Postgres...")
	// Fetch database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Build the connection string (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres %v", err)
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed opening connection to postgres %v", err)
	}

	// Set a global client instance
	Client = client

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	log.Println("Database connection established and migrations applied.")
}

// CloseDB closes the database connection.
func CloseDB() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
