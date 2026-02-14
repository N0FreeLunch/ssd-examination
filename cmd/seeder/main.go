package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"examination/cmd/seeder/internal/seeds"
	"examination/internal/ent"

	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func main() {
	// 1. Check environment (Prevent running in production)
	if os.Getenv("APP_ENV") == "production" {
		log.Fatal("Cannot run seeder in production environment")
	}

	// 2. Database Connection
	// Connect to SQLite file. Adjust path as needed.
	// In Docker, it's mapped to /data/local.db. Locally, it might be ./data/local.db
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "file:data/local.db?cache=shared&_pragma=foreign_keys(1)"
	} else {
		dbPath = fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", dbPath)
	}

	client, err := ent.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	// 3. Parse Flags
	seedName := flag.String("name", "exam_preview", "Name of the seeder to run (e.g., exam_preview)")
	clean := flag.Bool("clean", false, "Clean existing data before seeding")
	flag.Parse()

	ctx := context.Background()

	// 3.5 Auto-migrate Schema (Ensure tables exist)
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// 4. Run Seeder
	switch *seedName {
	case "exam_preview":
		if *clean {
			// Basic cleanup of exams logic could be here, or inside the seeder function
			log.Println("Cleaning up exam_preview data...")
			// Implementation detail: Delete previous exam with specific title
		}
		if err := seeds.SeedExamPreview(ctx, client); err != nil {
			log.Fatalf("Failed to seed exam preview: %v", err)
		}
		log.Println("Seeding exam_preview completed successfully.")
	default:
		log.Fatalf("Unknown seed name: %s", *seedName)
	}
}
