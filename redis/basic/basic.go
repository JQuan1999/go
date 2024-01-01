package basic

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "10.177.54.121:6379",
		Password: "123456",
		DB:       0,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	} else {
		fmt.Println("connect redis success")
	}
}

func Get() {
	val, err := rdb.Get(context.Background(), "key").Result()
	switch {
	case err == redis.Nil: // redis.Nil用来区分空的string回复和nil回复(表示key不存在)
		fmt.Println("key does not exist")
	case err != nil:
		fmt.Println("Get failed", err)
	case val == "":
		fmt.Println("value is empty")
	}
}

// 使用单个redis连接而不是连接池
func Conn() {
	conn := rdb.Conn()
	defer conn.Close()

	if err := conn.ClientSetName(context.Background(), "myclient").Err(); err != nil {
		panic(err)
	}
	name, err := conn.ClientGetName(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client name", name)
}
