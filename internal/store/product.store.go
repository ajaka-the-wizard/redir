package store

import (
	"context"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/models"
)

func (s *Store) CreatePrivateKey(ctx context.Context, productId int, hash string) (*models.Product, error) {
	return s.repo.CreatePrivateKey(ctx, productId, hash)
}

func (s *Store) GetProductById(ctx context.Context, productId int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	p, err := s.r.GetProduct(ctx, productId)
	if err == nil {
		return p, nil
	}
	p, err = s.repo.GetProductById(ctx, productId)
	if err != nil {
		return nil, err
	}
	s.r.SetProduct(ctx, *p)
	return p, nil
}

func (s *Store) CreateProduct(ctx context.Context, data *domain.CreateProductDetails) (*models.Product, error) {
	return s.repo.CreateProduct(ctx, data)
}
