package repository

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(pool *pgxpool.Pool, user *domain.
	CreateUserDetails, cfg *configs.EnvData) error {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXT_TIMEOUT)
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

func CreateOrLinkOauth(pool *pgxpool.Pool, cfg *configs.EnvData, id_or_sub string, email string, name string, provider string) (*domain.LightUser, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXT_TIMEOUT)
	defer cancel()
	var lUser domain.LightUser
	query := `
	INSERT INTO users (provider_sub, email, full_name, provider,verified)
	VALUES ($1, $2, $3, $4, true)
	ON CONFLICT (email) DO UPDATE SET
	verified = true
	provider_sub = EXCLUDED.provider_sub,
	provider = EXCLUDED.provider
	RETURNING id, email, admin, paid
	`
	err := pool.QueryRow(ctx, query, id_or_sub, email, name, provider).Scan(
		&lUser.Id,
		&lUser.Email,
		&lUser.Admin,
		&lUser.Paid,
	)
	if err != nil {
		return nil, err
	}
	return &lUser, nil
}

func GetUserByEmail(pool *pgxpool.Pool, cfg *configs.EnvData, email string) (*models.User, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXT_TIMEOUT)
	defer cancel()
	query := `
	SELECT id, email, password, verified
	FROM users
	WHERE email = $1
	`
	var user models.User
	err := pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Verified,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserById(pool *pgxpool.Pool, cfg *configs.EnvData, id uuid.UUID) (*models.User, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXT_TIMEOUT)
	defer cancel()
	query := `
	SELECT id, email, verified, full_name, paid
	FROM users
	WHERE id = $1
	`
	var user models.User
	err := pool.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Email,
		&user.Verified,
		&user.FullName,
		&user.Paid,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByProvider(pool *pgxpool.Pool, cfg *configs.EnvData, provider string, sub string) (*domain.LightUser, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXT_TIMEOUT)
	defer cancel()
	query := `
	SELECT id, email, paid, admin
	FROM users
	WHERE provider = $1 AND provider_sub = $2
	`
	var user domain.LightUser
	err := pool.QueryRow(ctx, query, provider, sub).Scan(
		&user.Id,
		&user.Email,
		&user.Paid,
		&user.Admin,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
