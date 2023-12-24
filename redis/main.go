package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var redisDb *redis.Client

func initClient() (err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err = redisDb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := initClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("连接成功")
}
