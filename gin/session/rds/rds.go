package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitCache(address string, password string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}

func GetRedis() *redis.Client {
	return redisClient
}
