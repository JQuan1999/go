package main

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type HistoryManager struct {
	// TODO:加锁或者使用sync.Map
	history  map[string]*SqlStats
	redisCli *redis.Client
}

// 找到对应key的历史记录sql stats
// 如果map中不存在, 从redis里面读
// redis中没读到说明是新的proxy
func (h *HistoryManager) find(key string) (*SqlStats, error) {
	if stats, ok := h.history[key]; ok {
		return stats, nil
	}
	res, err := h.redisCli.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	var sqlStats SqlStats
	stats := make([]SqlStatsRow, 0)
	if err := json.Unmarshal([]byte(res), &stats); err != nil {
		return nil, err
	}
	sqlStats.stats = stats
	return &sqlStats, nil
}

func (h *HistoryManager) set(key string, stats *SqlStats) {
	h.history[key] = stats
}
