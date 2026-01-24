package main

import (
	"database/sql"
	"examination/internal/api"
	"examination/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

func main() {
	// 1. Initialize DB
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "examination.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// 2. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 3. Health Check (DB Verification)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Try to write to the DB to ensure it's writable (and Litestream can replicate it)
		// We use a dedicated health check table.
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS health_check (id INTEGER PRIMARY KEY);")
		if err != nil {
			http.Error(w, fmt.Sprintf("Health Check Failed (Write): %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	svc := service.NewServer()
	apiHandler := api.NewStrictHandler(svc, nil)

	api.HandlerFromMux(apiHandler, r)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
