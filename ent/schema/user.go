package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"

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
			Sensitive(),
		field.String("first_name").
			NotEmpty(),
		field.String("last_name").
			NotEmpty(),
		field.Bool("is_admin").
			Default(false),
		field.UUID("uuid", uuid.UUID{}).
			Default(uuid.New),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}
