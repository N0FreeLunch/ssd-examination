package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"examination/internal/ent"
	"examination/internal/ent/exam"

	"entgo.io/ent/dialect"
	"modernc.org/sqlite"
)

func init() {
	sql.Register("sqlite3", &sqlite.Driver{})
}

func main() {
	// Create an in-memory sqlite database.
	// We use the same DSN settings as migration tool for consistency.
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	fmt.Println(">> Schema created successfully!")

	// 1. Create an Exam
	khExam, err := client.Exam.Create().
		SetTitle("Sungsil Univ 2024 Transfer Math").
		SetDescription("The 2024 transfer examination for Sungsil University.").
		SetTimeLimit(60).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating exam: %v", err)
	}
	fmt.Printf(">> Created Exam: %s (ID: %d)\n", khExam.Title, khExam.ID)

	// 2. Create Hierarchy (Section -> Topic -> Unit)
	// We can use the fluent builder to link them.

	// Section A
	sectA, err := client.Section.Create().
		SetTitle("Section A: Calculus").
		SetSeq(1).
		SetExam(khExam).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating section: %v", err)
	}

	// Topic 1
	top1, err := client.Topic.Create().
		SetTitle("Limits").
		SetSeq(1).
		SetExam(khExam).
		SetSection(sectA).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating topic: %v", err)
	}

	// Unit 1
	u1, err := client.Unit.Create().
		SetTitle("Limit Definition").
		SetSeq(1).
		SetExam(khExam).
		SetSection(sectA).
		SetTopic(top1).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating unit: %v", err)
	}

	fmt.Printf(">> Created Hierarchy: %s -> %s -> %s\n", sectA.Title, top1.Title, u1.Title)

	// 3. Query the data to verify relationships
	queryExam, err := client.Exam.Query().
		Where(exam.ID(khExam.ID)).
		WithSections(func(sq *ent.SectionQuery) {
			sq.WithTopics(func(tq *ent.TopicQuery) {
				tq.WithUnits()
			})
		}).
		Only(ctx)
	if err != nil {
		log.Fatalf("failed querying exam: %v", err)
	}

	fmt.Println("\n>> Query Result Tree:")
	fmt.Printf("Exam: %s\n", queryExam.Title)
	for _, s := range queryExam.Edges.Sections {
		fmt.Printf("  Section: %s\n", s.Title)
		for _, t := range s.Edges.Topics {
			fmt.Printf("    Topic: %s\n", t.Title)
			for _, u := range t.Edges.Units {
				fmt.Printf("      Unit: %s\n", u.Title)
			}
		}
	}
}
