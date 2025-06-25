package db

import (
	"context"
	"fmt"
	"os"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

func Init() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		return fmt.Errorf("no se pudo conectar a Redis: %v", err)
	}
	return nil
}
