package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	ctx := context.Background()
	// Initialize the database connection
	pool, err := db.InitDB(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	log.Println("Database connection established")
	defer pool.Close()

	db := db.NewPostgresDB(pool)
	server := api.NewServer(db)

	log.Println("API initialized")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(":"+port, server.InitRoutes())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
