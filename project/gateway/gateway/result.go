package main

import (
	"github.com/ngaut/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type CollectResulter interface {
	update(string, *HistoryManager)
	// 写入普罗米修斯Metric channel
	write(chan<- prometheus.Metric)
}

type ProxyStatus struct {
	status map[string]string
}

func NewProxyStatus() *ProxyStatus {
	var status ProxyStatus
	status.status = make(map[string]string)
	return &status
}

func (r *ProxyStatus) update(key string, h *HistoryManager) {

}

func (r *ProxyStatus) write(ch chan<- prometheus.Metric) {

}

type SqlStatsRow map[string]string

type SqlStats struct {
	stats []SqlStatsRow
}

func NewSqlStats() *SqlStats {
	var sqlstats SqlStats
	sqlstats.stats = make([]SqlStatsRow, 0)
	return &sqlstats
}

func (r *SqlStats) write(ch chan<- prometheus.Metric) {

}

func (r *SqlStats) append(row SqlStatsRow) {
	r.stats = append(r.stats, row)
}

func (r *SqlStats) update(key string, h *HistoryManager) {
	// 从history中获取key的历史记录
	history, err := h.find(key)

	if err != nil {
		if err != redis.Nil {
			log.Warnf("find key: %s history result failed, err: %v", key, err)
			return
		} else {
			// err == redis.Nil历史记录不存在为新纪录, 存入h中
			h.set(key, r)
			return
		}
	}

	// 历史记录存在修正历史记录和当前结果
	if r.stats == nil {
		// 1. 采集结果是否为nil, nil表示当前采集周期内没有, 修改历史记录
		for _, historyRow := range history.stats {
			historyRow["SingleCost"] = "0"
			historyRow["Count"] = "0"
			historyRow["PlanCost"] = "0"
			historyRow["DBCost"] = "0"
			historyRow["HandleRespCost"] = "0"
			historyRow["MinCost"] = "0"
			historyRow["MaxCost"] = "0"
			historyRow["AvgCost"] = "0"
			historyRow["TotalCost"] = "0"
			historyRow["StartCount"] = "0"
			historyRow["FinishCount"] = "0"
			historyRow["ConcurrentCount"] = "0"
			historyRow["MaxConcurrentCount"] = "0"
		}
		r.stats = history.stats
		return
	} else {
		// 2. 采集结果不为nil, 进行修正
		// 2.1 对于新的指纹sql, 设置历史sql stats记录为0

		// 2.2 修改"FirstSeen"、"Count"、"ErrorCount"累加
	}

}

type Result struct {
	status *ProxyStatus
	stats  *SqlStats
}

func NewResult(status *ProxyStatus, stats *SqlStats) *Result {
	return &Result{status: status, stats: stats}
}

func (r *Result) write(ch chan<- prometheus.Metric) {

}

func (r *Result) update(key string, h *HistoryManager) {

}
