package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Unit holds the schema definition for the Unit entity.
type Unit struct {
	ent.Schema
}

// Fields of the Unit.
func (Unit) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.Int("seq").Comment("Sequence order"),
		field.Int("exam_id"),
		field.Int("section_id").Optional().Nillable(),
		field.Int("topic_id").Optional().Nillable(),
	}
}

// Edges of the Unit.
func (Unit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("exam", Exam.Type).
			Ref("units").
			Field("exam_id").
			Unique().
			Required(),
		edge.From("section", Section.Type).
			Ref("units").
			Field("section_id").
			Unique(),
		edge.From("topic", Topic.Type).
			Ref("units").
			Field("topic_id").
			Unique(),
		edge.To("problems", Problem.Type),
	}
}

// Hooks of the Unit.
func (Unit) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				// Filter: Only apply on Create or Update operations
				if !m.Op().Is(ent.OpCreate | ent.OpUpdate | ent.OpUpdateOne) {
					return next.Mutate(ctx, m)
				}

				// Check if "topic_id" is being set.
				if v, exists := m.Field("topic_id"); exists && v != nil {
					// If topic_id is set to a non-nil value, we must clear section_id.
					if err := m.ClearField("section_id"); err != nil {
						// generic ClearField should be supported by ent.Mutation
						return nil, err
					}
				}
				return next.Mutate(ctx, m)
			})
		},
	}
}
