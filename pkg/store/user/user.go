package user

import "entgo.io/ent/entc/integration/ent"

type Store struct {
	db *ent.Client
}

func NewStore(db *ent.Client) *Store {
	return &Store{
		db: db,
	}
}
