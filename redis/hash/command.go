package hash

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

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

// redis hash操作主要有2-3个元素组成：
//   - key - redis key 唯一标识
//   - field - hash数据的字段名
//   - value - 值，有些操作不需要值

// -   HSet - 根据key和field字段设置，field字段的值
// -   HGet - 根据key和field字段，查询field字段的值
// -   HGetAll - 根据key查询所有字段和值
// -   HIncrBy - 根据key和field字段，累加数值。
// -   HKeys - 根据key返回所有字段名
// -   HLen - 根据key，查询hash的字段数量
// -   HMGet - 根据key和多个字段名，批量查询多个hash字段值
// -   HMSet - 根据key和多个字段名和字段值，批量设置hash字段值
// -   HSetNX - 如果field字段不存在，则设置hash字段值
// -   HDel - 根据key和字段名，删除hash字段，支持批量删除hash字段
// -   HExists - 检测hash字段名是否存在。

func TestHset() {
	rdb.HSet(ctx, "user", "key1", "value1")
	rdb.HSet(ctx, "user", "address", "127.0.0.1:8888")
	rdb.HSet(ctx, "user", "key3", "value3", "key4", "value4")
}

func TestHGet() {
	address, err := rdb.HGet(ctx, "user", "address").Result() // hegt获取"user"对象的key=address
	if err != nil {
		fmt.Println("hget failed, err: ", err)
	} else {
		fmt.Println("hget success, address: ", address)
	}
}

func TestGetAll() {
	user, err := rdb.HGetAll(ctx, "user").Result() // 获取user的所有key和value
	if err != nil {
		panic(err)
	}
	for key, value := range user {
		fmt.Println("key= ", key, " value= ", value)
	}
}
