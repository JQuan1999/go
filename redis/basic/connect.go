package basic

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func TestConnect() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.177.54.121:6379",
		Password: "123456",
		DB:       0,
	})

	value, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("连接redis出错, 出错信息: %v", err)
	}
	fmt.Println("connect success, return value= ", value)
}
