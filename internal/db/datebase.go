package db

import (
	"awesomeProject1/config"
	"awesomeProject1/ent"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var Client *ent.Client

func InitDB() {
	config.Log.Info("Attempting to connect to Postgres...")
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

	Client = client
	log.Println("Successfully connected to postgres")
	config.Log.Info("Successfully connected to Postgres!")
}
