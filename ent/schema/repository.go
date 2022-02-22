package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Repository holds the schema definition for the Repository entity.
type Repository struct {
	ent.Schema
}

// Fields of the Repository.
func (Repository) Fields() []ent.Field {
	return ent.Field{
		field.Uint64("repo_id").
			Unique(),
		field.String("name").
			NotEmpty().
			MaxLen(20),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Repository.
func (Repository) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("repositories").
			Unique(),
		Required(),
		Field("user_id"),
		edge.To("records", Record.Type),
	}
}
