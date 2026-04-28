package store

import (
	"context"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (s *Store) CreateUser(ctx context.Context, pool *pgxpool.Pool, u *domain.
	CreateUserDetails, cfg *configs.EnvData) error {
	return repository.CreateUser(ctx, pool, u, cfg)
}

func (s *Store) CreateOrLinkOauth(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, id_or_sub string, email string, name string, provider string) (*domain.LightUser, error) {
	return repository.CreateOrLinkOauth(ctx, pool, cfg, id_or_sub, email, name, provider)
}

func (s *Store) GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, email string) (*models.User, error) {
	const by string = "email"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	u, err := s.r.GetFullUser(ctx, email, by)
	if err == nil {
		return u, nil
	}
	u, err = repository.GetUserByEmail(ctx, pool, cfg, email)
	if err != nil {
		return nil, err
	}
	s.r.SetFullUser(ctx, u.Email, by, *u)
	return u, nil
}

func (s *Store) GetUserById(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, id uuid.UUID) (*models.User, error) {
	const by string = "id"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	u, err := s.r.GetFullUser(ctx, id.String(), by)
	if err == nil {
		return u, nil
	}
	u, err = repository.GetUserById(ctx, pool, cfg, id)
	if err != nil {
		return nil, err
	}
	s.r.SetFullUser(ctx, u.Id.String(), by, *u)
	return u, nil
}

func (s *Store) GetUserByProvider(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, provider string, sub string) (*domain.LightUser, error) {
	return repository.GetUserByProvider(ctx, pool, cfg, provider, sub)
}
