package pool

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Pool struct {
	rds     *redis.Client
	ctx     context.Context
	workers []*Worker
	ch      chan *Task
}

func NewPool(ctx context.Context, workerSize int, queueSize int) *Pool {
	var pool Pool
	pool.ctx = ctx
	pool.ch = make(chan *Task, queueSize)
	for i := 0; i < workerSize; i++ {
		worker := NewWorker(ctx, pool.ch)
		worker.Start()
		pool.workers = append(pool.workers, worker)
	}
	return &pool
}

func (p *Pool) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			close(p.ch)
			fmt.Println("Pool is closed")
		case <-ticker.C:
			now := time.Now().Unix()                                     // 获取时间戳
			p.rds.Set(p.ctx, "pool_key", fmt.Sprintf("pool_%d", now), 0) // keepalive
		}
	}
}

type Task struct {
	result chan int
	val    int
}

func (t *Task) Wait() int {
	return <-t.result
}

func (t *Task) Process() {
	t.result <- t.val * 2
}

type Worker struct {
	ch  <-chan *Task
	ctx context.Context
}

func NewWorker(ctx context.Context, ch <-chan *Task) *Worker {
	var w Worker
	w.ctx = ctx
	w.ch = ch
	return &w
}

func (w *Worker) Start() {
	for {
		select {
		case task := <-w.ch:
			if task == nil {
				return
			}
			task.Process()
		case <-w.ctx.Done():
			fmt.Println("worker is closed")
			return
		}
	}
}
