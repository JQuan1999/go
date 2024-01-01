package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ngaut/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

var sqlStatsProxyInstanceKey = "sql_stats_proxy"

type TaskType int

const (
	proxyInfo TaskType = 1
	sqlStats  TaskType = 2
	all       TaskType = 3
)

type ProxyInstance struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type Exporter struct {
	taskType TaskType
	redisCli *redis.Client
	workers  []*Worker
}

// 创建exporter
func NewExporter(taskType TaskType, redisCli *redis.Client) *Exporter {
	return &Exporter{
		redisCli: redisCli,
		taskType: taskType,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ctx, cfn := context.WithTimeout(context.Background(), time.Second*100)
	res, err := e.redisCli.Get(ctx, sqlStatsProxyInstanceKey).Result()
	if err != nil {
		log.Info("get proxy key from redis fail, error: %s\n", err)
	}
	cfn()
	var insts []*ProxyInstance
	if res != "" {
		json.Unmarshal([]byte(res), &insts)
		if len(insts) == 0 {
			log.Warn("no proxy insts")
			return
		}
	}
	for _, proxy := range insts {
		worker, err := NewWorker(e.taskType, proxy.Address, proxy.Port)
		if err != nil {
			log.Warnf("open proxy %s:%d failed, err:%v", proxy.Address, proxy.Port, err)
			continue
		}
		e.workers = append(e.workers, worker)
	}
	var wg sync.WaitGroup

	for _, worker := range e.workers {
		wg.Add(1)
		go worker.Collect(&wg)
	}
	wg.Wait()

}

func newMetricPullHandler(taskType TaskType, redisCli *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reg := prometheus.NewRegistry()
		reg.MustRegister(NewExporter(taskType, redisCli))

		gathers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			reg,
		}
		h := promhttp.HandlerFor(gathers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

// 运行exporter
func (e *Exporter) ServerHttp(redisCli *redis.Client) {
	http.Handle("/collect/sqlstats", newMetricPullHandler(sqlStats, redisCli))
}

type Worker struct {
	address string
	port    int
	task    TaskType
	conn    *sql.DB
	result  CollectResulter
	h       *HistoryManager
}

func NewWorker(task TaskType, address string, port int) (*Worker, error) {
	var worker Worker
	worker.address = address
	worker.port = port

	//TODO:dbstore管理mysql连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?multiStatements=true", "admin", "aNekZX9CWyve@RzQkY", address, port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(time.Minute * 10)
	worker.conn = db

	return &worker, nil
}

func (w *Worker) Collect(wg *sync.WaitGroup) {
	defer wg.Done()

	// 采集数据
	if w.task == proxyInfo {
		w.fetchProxyStatus()
	} else if w.task == sqlStats {
		w.fetchSqlStats()
	}

	// 修正数据
	w.result.update(w.key(), w.h)

	// 写入普罗米修斯
}

func (w *Worker) key() string {
	return fmt.Sprintf("proxy_%s_%d", w.address, w.port)
}

func (w *Worker) fetchAll() {
	status, err := w.collectProxyStatus()
	if err != nil {
		log.Warn("collect proxy %s:%d status failed, err: %v", w.address, w.port, err)
		return
	}
	stats, err := w.collectSqlStats()
	if err != nil {
		log.Warn("collect proxy %s:%d status failed, err: %v", w.address, w.port, err)
		return
	}
	w.result = NewResult(status, stats)
}

func (w *Worker) fetchProxyStatus() {
	status, err := w.collectProxyStatus()
	if err != nil {
		log.Warn("collect proxy %s:%d status failed, err: %v", w.address, w.port, err)
		return
	}
	w.result = status
}

func (w *Worker) fetchSqlStats() {
	stats, err := w.collectSqlStats()
	if err != nil {
		log.Warn("collect proxy %s:%d sql stats failed, err: %v", w.address, w.port, err)
		return
	}
	w.result = stats
}

func (w *Worker) collectProxyStatus() (*ProxyStatus, error) {
	type StringPair struct {
		key   string
		value string
	}
	rows, err := w.conn.Query("SHOW PROXY INFO")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proxyStatus ProxyStatus
	for rows.Next() {
		var pair StringPair
		if err = rows.Scan(&pair.key, &pair.value); err != nil {
			return nil, err
		}
		proxyStatus.status[pair.key] = pair.value
	}
	return &proxyStatus, err
}

func (w *Worker) collectSqlStats() (*SqlStats, error) {
	rows, err := w.conn.Query("show sql stats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats SqlStats
	for rows.Next() {
		desc := make(map[string]string)
		columns, err := rows.Columns()
		if err != nil {
			log.Errorf("common query get columns error: %v", err)
			return nil, err
		}
		args := make([]interface{}, 0)
		for _, column := range columns {
			desc[column] = ""
			args = append(args, new(string))
		}
		if err = rows.Scan(args...); err != nil {
			log.Warnf("scan failed, err: %v", err)
		}
		for i, column := range columns {
			desc[column] = *(args[i].(*string))
		}
		stats.append(desc)
	}
	return &stats, nil
}
