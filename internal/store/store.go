package store

import "github.com/ajaka-the-wizard/redir/internal/cache"

type Store struct {
	r *cache.Sredis
}

func InitializeStore(r *cache.Sredis) *Store {
	return &Store{
		r,
	}
}
