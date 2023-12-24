package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.177.54.121:6379",
		Password: "123456",
		DB:       0,
	})
	val, err := rdb.Get(context.Background(), "key").Result()
	fmt.Println(val)

}
