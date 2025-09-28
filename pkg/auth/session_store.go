package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jjnt224/chat8/pkg/config"
	"github.com/redis/go-redis/v9"
)

type SessionData struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

type SessionStore struct {
	Client *redis.Client
}

func NewSessionStore(cfg config.Config) *SessionStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   0,
	})
	return &SessionStore{Client: rdb}
}

func (s *SessionStore) Save(ctx context.Context, token string, data SessionData, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, "session:"+token, bytes, ttl).Err()
}

func (s *SessionStore) Get(ctx context.Context, token string) (*SessionData, error) {
	val, err := s.Client.Get(ctx, "session:"+token).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *SessionStore) Delete(ctx context.Context, token string) error {
	return s.Client.Del(ctx, "session:"+token).Err()
}
