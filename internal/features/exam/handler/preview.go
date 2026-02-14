package handler

import (
	"context"
	"examination/internal/ent"
	"examination/internal/ent/choice"
	"examination/internal/ent/exam"
	"examination/internal/ent/problemtranslation"
	"examination/internal/ent/section"
	"examination/internal/ent/unit"
	"html/template"
	"net/http"
	"path/filepath"
)

type ExamPreviewHandler struct {
	client *ent.Client
}

func NewExamPreviewHandler(client *ent.Client) *ExamPreviewHandler {
	return &ExamPreviewHandler{client: client}
}

func (h *ExamPreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Fetch the specific exam
	// Using the title from the seeder
	targetExam, err := h.client.Exam.Query().
		Where(exam.TitleEQ("Distributed Systems 101")).
		WithSections(func(sq *ent.SectionQuery) {
			sq.Order(ent.Asc(section.FieldSeq)).
				WithUnits(func(uq *ent.UnitQuery) {
					uq.Order(ent.Asc(unit.FieldSeq)).
						WithProblems(func(pq *ent.ProblemQuery) {
							pq.WithTranslations(func(ptq *ent.ProblemTranslationQuery) {
								ptq.Where(problemtranslation.LocaleEQ("en")).
									WithChoices(func(cq *ent.ChoiceQuery) {
										cq.Order(ent.Asc(choice.FieldSeq))
									})
							})
						})
				})
		}).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			http.Error(w, "Exam not found. Did you run the seeder?", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to load exam: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare Template
	// For simplicity, we parse the template on each request during dev.
	// In prod, this should be cached.
	tmplPath := filepath.Join("internal", "features", "exam", "ui", "exam_preview.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render
	if err := tmpl.Execute(w, targetExam); err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}
