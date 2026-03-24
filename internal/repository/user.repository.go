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
	CreateUserDetails, cfg *configs.EnvData) (*domain.LightUser, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXT_TIMEOUT)
	defer cancel()
	var lUser domain.LightUser
	query := `
	INSERT INTO users (full_name,email,password)
	VALUES ($1, $2, $3)
	RETURNING id, email, admin, paid
	`
	err := pool.QueryRow(ctx, query, user.FullName, user.Email, user.Password).Scan(
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
