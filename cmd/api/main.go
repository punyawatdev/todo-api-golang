package main

import (
	"context"
	"errors"
	"log"
	"net/http"
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

	// Setup Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connect to database
	dbConfig := database.LoadConfigFromEnv()
	db, err := database.NewPostgresDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		log.Println("Closing database connection...")
		db.Close()
	}()

	log.Println("✅ Connected to PostgreSQL")

	// Run Migrations
	if err := database.RunMigrations(ctx, db); err != nil {
    log.Fatalf("Migration failed: %v", err)
	}

	// Create & Start Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New(":"+port, db)

	// Run Server in Goroutine
	go func() {
		log.Printf("🚀 Server starting on :%s\n", port)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server forced to shutdown: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C) or SIGTERM (from Docker/K8s)
	<-ctx.Done()

	// Receiving the signal, begin the shutdown process.
	log.Println("Shutting down gracefully... (Press Ctrl+C again to force)")

	// Create a new context with timeout for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Attempt Graceful Shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Graceful shutdown failed: %v", err)
	}

	log.Println("👋 Server exited")
}
