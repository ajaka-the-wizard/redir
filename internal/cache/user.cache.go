package cache

import (
	"context"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/models"
)

func (r *Sredis) GetFullUser(ctx context.Context, identifier string, by string) (*models.User, error) {
	var user models.User
	key := "user-" + by + ":" + identifier
	err := r.rdb.HGetAll(ctx, key).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Sredis) SetFullUser(ctx context.Context, identifier string, by string, u models.User) error {
	key := "user-" + by + ":" + identifier
	exp := time.Minute * 5
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, u)
	pipe.Expire(ctx, key, exp)
	_, err := pipe.Exec(ctx)
	return err
}
