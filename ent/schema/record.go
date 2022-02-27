package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Record holds the schema definition for the Record entity.
type Record struct {
	ent.Schema
}

// Fields of the Record.
func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.Int("repo_id").
			Positive(),
		field.String("content").
			Annotations(entsql.Annotation{
				Size: 256,
			}).
			Validate(MaxRuneCount(256)),
		field.String("image_path"),
		//			Match(""), //TODO check prefix image url
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Record.
func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Ref("records").
			Unique().
			Required().
			Field("repo_id"),
	}
}
