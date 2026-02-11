package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo-api/internal/database"
	"todo-api/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	dbConfig := database.LoadConfigFromEnv()
	db, err := database.NewPostgresDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to PostgreSQL")

	// Create server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New(":"+port, db)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down gracefully...")
		cancel()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	log.Printf("🚀 Server starting on :%s\n", port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
