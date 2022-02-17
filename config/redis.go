package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	RedisDefault *redis.Client
	RedisCache   *redis.Client
}

func SetupRedis() (*Redis, error) {
	// Connect default redis
	dbDefault, err := strconv.Atoi(os.Getenv("REDIS_DEFAULT_DB"))
	redisDefault := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%s`, os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbDefault,
	})

	// Connect default cache
	redisCache := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%s`, os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	redis := Redis{
		RedisDefault: redisDefault,
		RedisCache:   redisCache,
	}

	return &redis, err
}
