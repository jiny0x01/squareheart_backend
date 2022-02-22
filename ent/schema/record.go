package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
)

// Record holds the schema definition for the Record entity.
type Record struct {
	ent.Schema
}

// Fields of the Record.
func (Record) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("record_id"),
		field.String("content").
			Annotations(entsql.Annotation{
				Size: 256
			}).
			Validate(MaxRuneCount(256)),
		field.String("image_path").
			Match(""), //TODO check prefix image url
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Record.
func (Record) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("repository", Repository.Type).
			Unique(),
			Required(),
			Field("repo_id"),
	}
}
