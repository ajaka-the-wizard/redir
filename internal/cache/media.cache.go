package cache

import (
	"context"

	"github.com/ajaka-the-wizard/redir/internal/models"
	"github.com/redis/go-redis/v9"
)

func (r *Sredis) GetMedia(ctx context.Context, publicKey string) (*models.Media, error) {
	var media models.Media
	key := "f" + publicKey
	res, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, redis.Nil
	}
	err = r.rdb.HGetAll(ctx, key).Scan(&media)
	return &media, nil
}
