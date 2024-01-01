package history

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var address string
var auth string
var ctx context.Context

type FullResult struct {
	Digest  string
	SQLStat map[string]string
}

func init() {
	address = "rediscache2.cdb.17usoft.com:3630"
	auth = "tcbase.dbmiddleware.dbproxy.scanner:qa:inst_data@TCBase.Cache.v3:b140320d"
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: auth,
	})
	if result := rdb.Ping(ctx); result.Err() != nil {
		panic(result.Err())
	}
}

func GetProxyAddress() {

}

func TestPipelineGet() {
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	} else {
		fmt.Println("ping success")
	}
	// TODO:get proxyinfo from redis
	proxyAddress, err := dbstore.GetProxyInst() // get proxy info from mysql
	if err != nil {
		panic(err)
	}
	fmt.Println("find proxy count: ", len(proxyAddress))
	// dataBase := []string{"test_ms_proxy"}

	proxyKeys := make([]string, 0)
	// 创建redis-pipeline
	pipeClient := rdb.Pipeline()
	for i := 0; i < len(proxyAddress); i++ {
		key := GetProxyDigestKey(proxyAddress[i])
		proxyKeys = append(proxyKeys, key)
		pipeClient.Get(ctx, key) // get命令加入pipeline
	}

	startTime := time.Now()
	res, err := pipeClient.Exec(ctx) // send all the commands buffered in the pipeline to the redis-server
	endTime := time.Since(startTime) // test fetch data cost time
	fmt.Printf("fetch history data from redis cost: %f second\n", endTime.Seconds())

	if err != nil && err != redis.Nil {
		fmt.Printf("pipeline execute failed, kvs: %v, error: %v", proxyKeys, err)
		return
	}

	// 遍历pipeline的结果
	// 统计byte数量
	compressByteCount := 0
	normalByteCount := 0
	proxyCount := 0

	for idx, cmdRes := range res {

		cmd, ok := cmdRes.(*redis.StringCmd) // 转换为StringCmd
		if !ok {
			continue
		}
		historyData, err := cmd.Result()
		compressByteCount += len([]byte(historyData))

		switch {
		case err == redis.Nil:
			fmt.Printf("proxy: %s don't have history result\n", proxyAddress[idx])
		case err != nil:
			fmt.Printf("get proxy: %s history from redis fail, err: %v\n", proxyAddress[idx], err)
		default:
			proxyCount++
			data, err := DecodeByGzip(historyData) // 解压历史数据
			if err != nil {
				fmt.Println("decode failed, err: ", err)
				continue
			}
			count := len([]byte(data)) // 占用的byte数
			normalByteCount += count

			var res []FullResult
			if err := json.Unmarshal(data, &res); err != nil { // 反序列化历史数据
				fmt.Printf("unmarsal failed, err: %v\n", err)
				continue
			}

			fmt.Printf("proxy: %s cost memory: %dkb, digest code count: %d\n", proxyAddress[idx], int64(count/1024), len(res))

			// for idx := range res {
			// 	fmt.Println("degest code = ", res[idx].Digest) // 打印指纹
			// 	for k, v := range res[idx].SQLStat {           // 打印SQLStat
			// 		fmt.Println(k, " = ", v)
			// 	}
			// }
		}
	}
	fmt.Printf("in fact proxy count= %d, history data compress cost mem= %d KB, normal mem= %d KB\n", proxyCount, int64(compressByteCount/1024), int64(normalByteCount/1024))
}

func TestPipelineWrite() {

}
