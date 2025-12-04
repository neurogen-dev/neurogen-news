package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(redisURL string) (*RedisClient, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// Cache helpers
func (r *RedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, data, expiration).Err()
}

func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Rate limiting
func (r *RedisClient) CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	count, err := r.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.Expire(ctx, key, window)
	}

	return count <= int64(limit), nil
}

// Session management
func (r *RedisClient) SetSession(ctx context.Context, sessionID string, userID string, expiration time.Duration) error {
	return r.Set(ctx, "session:"+sessionID, userID, expiration).Err()
}

func (r *RedisClient) GetSession(ctx context.Context, sessionID string) (string, error) {
	return r.Get(ctx, "session:"+sessionID).Result()
}

func (r *RedisClient) DeleteSession(ctx context.Context, sessionID string) error {
	return r.Del(ctx, "session:"+sessionID).Err()
}

// Pub/Sub for real-time notifications
func (r *RedisClient) PublishNotification(ctx context.Context, userID string, notification interface{}) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	return r.Publish(ctx, "notifications:"+userID, data).Err()
}

func (r *RedisClient) SubscribeNotifications(ctx context.Context, userID string) *redis.PubSub {
	return r.Subscribe(ctx, "notifications:"+userID)
}

