package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name"),
		field.String("email").Unique(),
		field.Int("age"),
		field.String("pwd").
			Sensitive().
			NotEmpty().
			MinLen(8).
			MaxLen(128),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
