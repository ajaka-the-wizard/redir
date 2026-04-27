package memory

import (
	"context"
	"log"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	"github.com/ajaka-the-wizard/redir/internal/domain"
	"github.com/redis/go-redis/v9"
)

type Sredis struct {
	rdb *redis.Client
}

func InitializeRedis(ctx context.Context, cfg *configs.EnvData) *Sredis {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_ADDR,
		Password: cfg.REDIS_PASSWORD,
		DB:       0,
		Protocol: 2,
	})
	err := rdb.Set(ctx, "ping", "pong", 0).Err()
	if err != nil {
		panic(err)
	}
	log.Println("Redis connected successfully")
	return &Sredis{
		rdb,
	}
}

func (r *Sredis) Clean() {
	r.rdb.Close()
}

func (r *Sredis) SetUserOnline(ctx context.Context, sessionId string, u *domain.
	LightUser) (time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	key := "user:" + sessionId
	exp := time.Hour * 24
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, *u)
	pipe.Expire(ctx, key, exp)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().Add(exp), nil
}

func (r *Sredis) GetUser(ctx context.Context, sessionId string) (*domain.
	LightUser, bool) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var user domain.LightUser
	key := "user:" + sessionId
	err := r.rdb.HGetAll(ctx, key).Scan(&user)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func (r *Sredis) RevokeUser(ctx context.Context, sessionId string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	key := "user:" + sessionId
	err := r.rdb.Del(ctx, key).Err()
	return err
}
