package main

import (
	"context"
	"database/sql"
	"examination/internal/ent"
	"examination/internal/features/exam/handler"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func main() {
	// 1. Initialize DB (Ent Client)
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Default to local development path inside data directory
		dbPath = "file:data/local.db?cache=shared&_pragma=foreign_keys(1)"
	}

	client, err := ent.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer client.Close()

	// 2. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 3. Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Validating DB connection by running a simple query
		if _, err := client.Exam.Query().Limit(1).Count(r.Context()); err != nil {
			http.Error(w, fmt.Sprintf("Health Check Failed: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	var wg sync.WaitGroup

	// 4. Feature Handlers
	examHandler := handler.NewExamPreviewHandler(client)
	r.Get("/exams/preview", examHandler.ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Examination Service - SSR/HTMX Mode"))
	})

	// Temporary route for Graceful Shutdown testing (Goroutine scenario)
	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request on /slow. Spawning background task...")

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("Background task: Starting 5 second job...")
			time.Sleep(5 * time.Second)
			log.Println("Background task: Finished 5 second job.")
		}()

		w.Write([]byte("Request received. Background task started!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8180" // Default changed to 8180 to avoid conflicts
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Signal received, starting graceful shutdown...")

	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Wait for all background goroutines to finish
	log.Println("Waiting for background tasks to complete...")
	wg.Wait()

	log.Println("Server exiting")
}
