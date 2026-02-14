package seeds

import (
	"context"
	"examination/internal/ent"
	"examination/internal/ent/choice"
	"examination/internal/ent/exam"
	"examination/internal/ent/problem"
	"examination/internal/ent/problemtranslation"
	"examination/internal/ent/section"
	"examination/internal/ent/unit"
	"fmt"
	"log"
	"time"
)

// SeedExamPreview seeds data for the Exam Preview prototype scenario.
func SeedExamPreview(ctx context.Context, client *ent.Client) error {
	const examTitle = "Distributed Systems 101"

	// 1. Clean up existing exam data to avoid duplicates (simplified cleanup)
	// Find existing exam by title
	// Note: In a real seeder, we might want to be more careful.
	// For now, if it exists, we delete it and recreate.
	existingExam, err := client.Exam.Query().
		Where(exam.TitleEQ(examTitle)).
		Only(ctx)

	if err == nil {
		log.Printf("Deleting existing exam: %s", existingExam.Title)

		// Manual Cascade Delete (Bottom-Up)
		// 1. Choices
		_, err := client.Choice.Delete().Where(
			choice.HasProblemTranslationWith(
				problemtranslation.HasProblemWith(
					problem.HasUnitWith(
						unit.HasExamWith(exam.ID(existingExam.ID)),
					),
				),
			),
		).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed deleting choices: %w", err)
		}

		// 2. Translations
		_, err = client.ProblemTranslation.Delete().Where(
			problemtranslation.HasProblemWith(
				problem.HasUnitWith(
					unit.HasExamWith(exam.ID(existingExam.ID)),
				),
			),
		).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed deleting translations: %w", err)
		}

		// 3. Problems
		_, err = client.Problem.Delete().Where(
			problem.HasUnitWith(
				unit.HasExamWith(exam.ID(existingExam.ID)),
			),
		).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed deleting problems: %w", err)
		}

		// 4. Units
		_, err = client.Unit.Delete().Where(
			unit.HasExamWith(exam.ID(existingExam.ID)),
		).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed deleting units: %w", err)
		}

		// 5. Sections
		_, err = client.Section.Delete().Where(
			section.HasExamWith(exam.ID(existingExam.ID)),
		).Exec(ctx)
		if err != nil {
			return fmt.Errorf("failed deleting sections: %w", err)
		}

		// 6. Exam
		if err := client.Exam.DeleteOne(existingExam).Exec(ctx); err != nil {
			return fmt.Errorf("failed deleting existing exam: %w", err)
		}
	} else if !ent.IsNotFound(err) {
		return fmt.Errorf("failed querying existing exam: %w", err)
	}

	// 2. Create Exam
	exam, err := client.Exam.Create().
		SetTitle(examTitle).
		SetDescription("An introductory exam covering fundamental concepts of distributed systems.").
		SetTimeLimit(60).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating exam: %w", err)
	}
	log.Printf("Created Exam: %s (ID: %d)", exam.Title, exam.ID)

	// 3. Create Section
	section, err := client.Section.Create().
		SetTitle("Data Consistency").
		SetSeq(1).
		SetExam(exam). // Use ID or Edge
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating section: %w", err)
	}
	log.Printf("Created Section: %s", section.Title)

	// 4. Create Units (Problems)
	// Unit 1: Simple Multiple Choice
	u1, err := client.Unit.Create().
		SetTitle("CAP Theorem Basics").
		SetSeq(1).
		SetExam(exam).
		SetSection(section).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating unit 1: %w", err)
	}

	p1, err := client.Problem.Create().
		SetUnit(u1).
		SetType("SOURCE").
		SetDifficulty(1).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating problem 1: %w", err)
	}

	pt1, err := client.ProblemTranslation.Create().
		SetProblem(p1).
		SetLocale("en").
		SetTitle("CAP Theorem").
		SetContent("In the CAP theorem, which two properties cannot be simultaneously guaranteed in a distributed system with network partitions?").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating translation 1: %w", err)
	}

	_, err = client.Choice.CreateBulk(
		client.Choice.Create().SetProblemTranslation(pt1).SetContent("Consistency & Availability").SetIsCorrect(true).SetSeq(1),
		client.Choice.Create().SetProblemTranslation(pt1).SetContent("Availability & Partition Tolerance").SetIsCorrect(false).SetSeq(2),
		client.Choice.Create().SetProblemTranslation(pt1).SetContent("Consistency & Partition Tolerance").SetIsCorrect(false).SetSeq(3),
		client.Choice.Create().SetProblemTranslation(pt1).SetContent("Reliability & Scalability").SetIsCorrect(false).SetSeq(4),
	).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating choices 1: %w", err)
	}

	// Unit 2: Markdown Content (Eventual Consistency)
	u2, err := client.Unit.Create().
		SetTitle("Eventual Consistency").
		SetSeq(2).
		SetExam(exam).
		SetSection(section).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating unit 2: %w", err)
	}

	p2, err := client.Problem.Create().
		SetUnit(u2).
		SetType("SOURCE").
		SetDifficulty(2).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating problem 2: %w", err)
	}

	const mdContent = `
### Understanding Eventual Consistency

Eventual consistency is a consistency model used in distributed computing to achieve high availability. It guarantees that, if no new updates are made to a given data item, eventually all accesses to that item will return the last updated value.

**Which of the following statements is true regarding Eventual Consistency?**
`

	pt2, err := client.ProblemTranslation.Create().
		SetProblem(p2).
		SetLocale("en").
		SetTitle("Eventual Consistency Details").
		SetContent(mdContent).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating translation 2: %w", err)
	}

	_, err = client.Choice.CreateBulk(
		client.Choice.Create().SetProblemTranslation(pt2).SetContent("Data is instantly replicated to all nodes.").SetIsCorrect(false).SetSeq(1),
		client.Choice.Create().SetProblemTranslation(pt2).SetContent("It allows for temporary inconsistencies but converges over time.").SetIsCorrect(true).SetSeq(2),
	).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating choices 2: %w", err)
	}

	// Unit 3: Code Block (Go Channel)
	u3, err := client.Unit.Create().
		SetTitle("Go Channel Behavior").
		SetSeq(3).
		SetExam(exam).
		SetSection(section).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating unit 3: %w", err)
	}

	p3, err := client.Problem.Create().
		SetUnit(u3).
		SetType("SOURCE").
		SetDifficulty(3).
		SetCreatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating problem 3: %w", err)
	}

	const codeContent = "What is the output of the following Go code?\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    ch := make(chan int, 1)\n    ch <- 1\n    fmt.Println(<-ch)\n}\n```"

	pt3, err := client.ProblemTranslation.Create().
		SetProblem(p3).
		SetLocale("en").
		SetTitle("Go Channels").
		SetContent(codeContent).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating translation 3: %w", err)
	}

	_, err = client.Choice.CreateBulk(
		client.Choice.Create().SetProblemTranslation(pt3).SetContent("1").SetIsCorrect(true).SetSeq(1),
		client.Choice.Create().SetProblemTranslation(pt3).SetContent("Deadlock").SetIsCorrect(false).SetSeq(2),
		client.Choice.Create().SetProblemTranslation(pt3).SetContent("Runtime Error").SetIsCorrect(false).SetSeq(3),
	).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating choices 3: %w", err)
	}

	return nil
}
