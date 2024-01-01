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

var clients []*redis.Client

// 在加锁逻辑里，我们主要是对每个Redis实例执行SET resource_name my_random_value NX PX 30000获取锁，
// 然后把成功获取锁的客户端放到一个channel里（这里使用slice可能有并发问题），同时使用sync.WaitGroup等待所有获取锁操作结束。
// 然后添加defer释放锁逻辑，释放锁逻辑很简单，只是把成功拿到的锁给释放掉即可。
// 最后判断成功获取到的锁的数量是否大于一半，如果没有得到一半以上的锁，说明加锁失败。
// 如果加锁成功接下来就是进行业务处理。
func multiLottery(ctx context.Context) error {
	// 加锁
	myRandomValue := gofakeit.UUID()
	resourceName := "resource_name"
	var wg sync.WaitGroup
	wg.Add(len(clients))
	// 这里主要是确保不要加锁太久，这样会导致业务处理的时间变少
	lockCtx, _ := context.WithTimeout(ctx, time.Millisecond*5)
	// 成功获得锁的Redis实例的客户端
	successClients := make(chan *redis.Client, len(clients))
	for _, client := range clients {
		go func(client *redis.Client) {
			defer wg.Done()
			ok, err := client.SetNX(lockCtx, resourceName, myRandomValue, time.Second*30).Result()
			if err != nil {
				return
			}
			if !ok {
				return
			}
			successClients <- client
		}(client)
	}
	wg.Wait() // 等待所有获取锁操作完成
	close(successClients)
	// 解锁，不管加锁是否成功，最后都要把已经获得的锁给释放掉
	defer func() {
		script := redis.NewScript(unlockScript)
		for client := range successClients {
			go func(client *redis.Client) {
				script.Run(ctx, client, []string{resourceName}, myRandomValue)
			}(client)
		}
	}()
	// 如果成功加锁得客户端少于客户端数量的一半+1，表示加锁失败
	if len(successClients) < len(clients)/2+1 {
		return errors.New("系统繁忙，请重试")
	}

	// 业务处理
	time.Sleep(time.Second)
	return nil
}

func TestMultiLock() {
	clients = append(clients, redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	}), redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   1,
	}), redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   2,
	}), redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   3,
	}), redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   4,
	}))
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err := multiLottery(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
		err := multiLottery(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
	time.Sleep(time.Second)
}
