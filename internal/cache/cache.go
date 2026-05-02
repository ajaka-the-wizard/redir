package cache

import (
	"context"
	"log"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/configs"
	fts "github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Sredis struct {
	rdb *redis.Client
}

func InitializeRedis(ctx context.Context, cfg *configs.EnvData) *Sredis {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

func structToInterface(s any) map[string]any {
	t := fts.New(s)
	t.TagName = "redis"
	m := t.Map()
	for k, v := range m {
		if u, ok := v.(uuid.UUID); ok {
			m[k] = u.String()
		}
	}
	return m
}
