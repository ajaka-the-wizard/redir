package repository

import (
	"context"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, pool *pgxpool.Pool, user *domain.
	CreateUserDetails, cfg *configs.EnvData) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	INSERT INTO users (full_name,email,password)
	VALUES ($1, $2, $3)
	`
	_, err := pool.Exec(ctx, query, user.FullName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func CreateOrLinkOauth(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, id_or_sub string, email string, name string, provider string) (*domain.LightUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	INSERT INTO users (provider_sub, email, full_name, provider,verified)
	VALUES ($1, $2, $3, $4, true)
	ON CONFLICT (email) DO UPDATE SET
	verified = true,
	provider_sub = EXCLUDED.provider_sub,
	provider = EXCLUDED.provider
	RETURNING id, email, admin, paid
	`
	rows, err := pool.Query(ctx, query, id_or_sub, email, name, provider)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.LightUser])
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	SELECT id, email, password, verified
	FROM users
	WHERE email = $1
	`
	rows, err := pool.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserById(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, id uuid.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	SELECT id, email, verified, full_name, paid
	FROM users
	WHERE id = $1
	`
	rows, err := pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByProvider(ctx context.Context, pool *pgxpool.Pool, cfg *configs.EnvData, provider string, sub string) (*domain.LightUser, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
	SELECT id, email, paid, admin
	FROM users
	WHERE provider = $1 AND provider_sub = $2
	`
	rows, err := pool.Query(ctx, query, provider, sub)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.LightUser])

	if err != nil {
		return nil, err
	}

	return &user, nil
}
