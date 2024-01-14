package listcommand

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

var ctx = context.Background()

// 链表支持的操作：LPush、RPush、LPop、RPop、LRange、LIndex
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

func Clear() {
	for {
		// 从链表右侧循环弹出链表元素
		val, err := rdb.RPop(ctx, "list").Result()
		if err == redis.Nil {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Println("val= ", val) // 依次打印1、2、3
	}
}

func TestLPush() {
	// 返回值是当前元素的数量
	// LPush链表左侧插入元素
	n, err := rdb.LPush(ctx, "list", 1, 2, 3).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(n)
	Clear()
}

func TestLIndex() {
	// 插入链表元素
	for i := 0; i < 10; i++ {
		err := rdb.LPush(ctx, "list", i).Err()
		if err != nil {
			panic(err)
		}
	}
	// LIndex返回链表下表对应的元素
	for i := 0; i < 10; i++ {
		val, err := rdb.LIndex(ctx, "list", int64(i)).Result()
		if err != nil {
			fmt.Println("get lindex element failed, err: ", err)
		} else {
			fmt.Println("i = ", i, " val= ", val) // val = 9, 8, 7, ..., 0
		}
	}
	// 返回链表的长度
	length, err := rdb.LLen(ctx, "list").Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("len=", length)
	}
	Clear()
}

func TestLRange() {
	// 测试LRange获取某个范围内的元素集
	err := rdb.LPush(ctx, "list", 1, 2, 3, 4, 5, 6, 7, 8).Err()
	if err != nil {
		panic(err)
	}
	// 获取0-5下标的元素
	vals, err := rdb.LRange(ctx, "list", 0, 5).Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("val=", vals)
	}
	// 获取0-100下标的元素
	vals, err = rdb.LRange(ctx, "list", 0, 100).Result() // 超出范围只会返回1,2,3,4,5,6,7,8
	if err != nil {
		panic(err)
	} else {
		fmt.Println("val=", vals)
	}
	Clear()
}
