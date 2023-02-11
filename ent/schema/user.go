package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"

	"github.com/MatthewFrisby/thesis-pieces/ent/schema/mixins"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			NotEmpty().
			Unique(),
		field.String("email").
			NotEmpty().
			Unique(),
		field.String("password").
			NotEmpty().
			MinLen(8).
			Sensitive(),
		field.Bool("is_admin").
			Default(false),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}
