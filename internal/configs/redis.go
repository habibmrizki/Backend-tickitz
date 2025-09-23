package configs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	rdbHost := os.Getenv("RDBHOST")
	rdbpPort := os.Getenv("RDBPORT")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", rdbHost, rdbpPort),
	})
	return rdb, nil
}

// Fungsi untuk blacklist token
func BlacklistToken(rdb *redis.Client, token string, duration time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, "blacklist:"+token, "true", duration).Err()
}

// Fungsi untuk cek token di blacklist
func IsTokenBlacklisted(rdb *redis.Client, token string) (bool, error) {
	ctx := context.Background()
	val, err := rdb.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "true", nil
}
