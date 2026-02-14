package main

import (
	"database/sql"
	"examination/internal/ent"
	"examination/internal/features/exam/handler"
	"fmt"
	"log"
	"net/http"
	"os"

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

	// 4. Feature Handlers
	examHandler := handler.NewExamPreviewHandler(client)
	r.Get("/exams/preview", examHandler.ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Examination Service - SSR/HTMX Mode"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8180" // Default changed to 8180 to avoid conflicts
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
