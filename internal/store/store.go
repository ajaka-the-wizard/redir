package store

import (
	"github.com/ajaka-the-wizard/redir/internal/cache"
	"github.com/ajaka-the-wizard/redir/internal/repository"
)

type Store struct {
	r    *cache.Sredis
	repo *repository.Repository
}

func InitializeStore(r *cache.Sredis, repo *repository.Repository) *Store {
	return &Store{
		r,
		repo,
	}
}
