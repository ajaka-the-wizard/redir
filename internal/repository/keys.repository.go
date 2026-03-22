package repository

import (
	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/ajaka-the-wizard/redir/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePrivateKey(pool *pgxpool.Pool, cfg *configs.EnvData, user uuid.UUID, hash string) (*models.ClientKeys, error) {
	ctx, cancel := utils.CreateContextWithStatedTime(cfg.CONTEXTTIMEOUT)
	defer cancel()
	var key models.ClientKeys
	query := `
	INSERT INTO client_keys (private_key, user_id)
	VALUES($1, $2)
	RETURNING client_id, created_at
	`
	err := pool.QueryRow(ctx, query, hash, user).Scan(
		&key.ClientId,
		&key.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &key, nil
}
