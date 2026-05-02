package cache

import (
	"context"
	"time"

	"github.com/ajaka-the-wizard/redir/internal/domain"
)

func (r *Sredis) Clean() {
	r.rdb.Close()
}

func (r *Sredis) SetUserOnline(ctx context.Context, sessionId string, u *domain.
	LightUser) (time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	key := "user:" + sessionId
	exp := time.Hour * 24
	m := structToInterface(*u)
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, m)
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
	s := r.rdb.HGetAll(ctx, key)
	u, err := s.Result()
	if err != nil || len(u) == 0 {
		return nil, false
	}
	err = s.Scan(&user)
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
