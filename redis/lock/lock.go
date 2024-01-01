package lock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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

const unlockScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.Call("del", KEYS[1])
else
	return 0
end`

// 1. 生成随机值
// 2. 使用SET resource_name my_random_value NX PX 30000加锁
// 3. 如果加锁失败，直接返回
// 4. defer添加解锁逻辑，保证在函数退出的时候会执行
// 5. 执行业务逻辑
func lottery(ctx context.Context) error {
	// 加锁
	myRandomValue := gofakeit.UUID()
	resourceName := "resource_name"
	ok, err := rdb.SetNX(ctx, resourceName, myRandomValue, time.Second*30).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("系统繁忙请重试")
	}
	// 解锁
	defer func() {
		script := redis.NewScript(unlockScript)
		script.Run(ctx, rdb, []string{resourceName}, myRandomValue)
	}()

	// 业务逻辑
	time.Sleep(time.Second)
	return nil
}

func Testlock() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err := lottery(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err := lottery(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
}
