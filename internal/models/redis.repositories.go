// redis.repositories.go
package models

// import (
// 	"context"
// 	"time"

// 	"github.com/redis/go-redis/v9"
// )

// type RedisRepository struct {
// 	client *redis.Client
// }

// func NewRedisRepository(client *redis.Client) *RedisRepository {
// 	return &RedisRepository{client: client}
// }

// // BlacklistToken adds a token to the blacklist with an expiration time.
// func (r *RedisRepository) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
// 	return r.client.Set(ctx, token, "blacklisted", expiration).Err()
// }

// // IsTokenBlacklisted checks if a token exists in the blacklist.
// func (r *RedisRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
// 	val, err := r.client.Get(ctx, token).Result()
// 	if err == redis.Nil {
// 		return false, nil
// 	}
// 	if err != nil {
// 		return false, err
// 	}
// 	return val == "blacklisted", nil
// }
