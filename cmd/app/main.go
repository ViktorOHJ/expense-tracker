package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/auth"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	ctx := context.Background()

	// Инициализация базы данных
	pool, err := db.InitDB(ctx, os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer pool.Close()

	database := db.NewPostgresDB(pool)

	// Инициализация сервисов авторизации
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"))
	passwordService := auth.NewPasswordService()

	server := api.NewServer(database, jwtService, passwordService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	err = http.ListenAndServe(":"+port, server.InitRoutes())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
