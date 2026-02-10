package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"examination/internal/ent/migrate"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func main() {
	ctx := context.Background()

	// Create a local migration directory.
	dir, err := atlas.NewLocalDir("migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}

	// Migrate options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.SQLite),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: go run cmd/migration/main.go <name>")
	}

	// Generate migrations using NamedDiff.
	// Atlas requires a "Dev Database" to calculate the diff.
	// We use an in-memory database with _pragma to enable FKs.
	// 'sqlite://' scheme is required by Atlas.
	err = migrate.NamedDiff(ctx, "sqlite://ent?mode=memory&cache=shared&_pragma=foreign_keys(1)", os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
