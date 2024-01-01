package string_command

import (
	"context"
	"fmt"
	"time"

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

func TestSet() {
	//Set方法的最后一个参数表示过期时间，0表示永不过期
	err := rdb.Set(context.Background(), "key1", "value1", 0).Err()
	if err != nil {
		panic(err)
	}

	// key2将在1分钟后过期
	err = rdb.Set(context.Background(), "key2", "value2", time.Minute*1).Err()
	if err != nil {
		panic(err)
	}
}

// setnx()和setex
// SexNX()仅当key不存在的时候才设置，如果key已经存在则不做任何操作，而SetEX()方法不管该key是否已经存在缓存中直接覆盖
func TestSetNx() {
	res1, err := rdb.SetNX(context.Background(), "key1", "value1", 0).Result() // result返回bool,err
	if err != nil {
		fmt.Println("SetNX failed, err: ", err)
	} else {
		fmt.Println("SetNX success, res: ", res1) // 返回false
	}
	res2, err := rdb.SetEx(context.Background(), "key1", "value1", time.Second*10).Result()
	if err != nil {
		fmt.Println("SetEx failed, err: ", err)
	} else {
		fmt.Println("SetEx success, res: ", res2)
	}
}

// Append()表示往字符串后面追加元素，返回值是字符串的总长度
func TestAppend() {
	ctx := context.Background()
	err := rdb.Set(ctx, "key", "hello", 0).Err() // set key
	if err != nil {
		panic(err)
	}
	length, err := rdb.Append(ctx, "key", " world!").Result() // append key
	if err != nil {
		panic(err)
	}
	fmt.Printf("当前缓存key的长度为: %v\n", length)  //12
	val, err := rdb.Get(ctx, "key").Result() // get key
	if err != nil {
		panic(err)
	}
	fmt.Printf("当前缓存key的值为: %v\n", val) //hello world!
}
