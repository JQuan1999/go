package history

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis() {
	var address = "rediscache2.cdb.17usoft.com:3630"
	var auth = "tcbase.dbmiddleware.dbproxy.scanner:qa:inst_data@TCBase.Cache.v3:b140320d"
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: auth,
	})
	if result := rdb.Ping(context.Background()); result.Err() != nil {
		panic(result.Err())
	}
}

func GetRedisClient() *redis.Client {
	// 加锁实现线程安全
	if rdb == nil {
		InitRedis()
	}
	return rdb
}
