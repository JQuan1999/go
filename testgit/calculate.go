package main

import (
	"context"
	"log"
	"sync"
)

type Calculator struct {
	workers   []Worker
	workerNum int
	taskQueue chan Tasker
	ctx       context.Context
}

func NewCalculator(ctx context.Context, workerNum, queueSize int) *Calculator {
	var cal Calculator
	cal.ctx = ctx
	cal.workerNum = workerNum
	cal.taskQueue = make(chan Tasker, queueSize)
	cal.workers = make([]Worker, workerNum)
	return &cal
}

func (cal *Calculator) Run() {
	var wg sync.WaitGroup
	wg.Add(cal.workerNum)

	for i := 0; i < cal.workerNum; i++ {
		cal.workers[i] = Worker{
			id: i,
			ch: cal.taskQueue,
			wg: &wg,
		}
		cal.workers[i].start()
	}
	<-cal.ctx.Done()
	wg.Wait()
}

func (cal *Calculator) Submit(task Tasker) {
	cal.taskQueue <- task
}

type Worker struct {
	id int
	wg *sync.WaitGroup
	ch <-chan Tasker
}

func (work *Worker) start() {
	log.Printf("worker [%d] start\n", work.id)
	quitFunc := func() {
		work.wg.Done()
		log.Printf("worker [%d] stop\n", work.id)
	}
	go func() {
		for task := range work.ch {
			task.Process()
		}
		quitFunc()
	}()
}

type Tasker interface {
	Process()
	Wait() int
}

type AddTask struct {
	Num1   int `form:"num1" binding:"required"`
	Num2   int `form:"num1" binding:"required"`
	result chan int
}

func (task *AddTask) Process() {
	task.result <- task.Num1 + task.Num2
}

func (task *AddTask) Wait() int {
	r := <-task.result
	return r
}

type SubTask struct {
	Num1   int `form:"num1" binding:"required"`
	Num2   int `form:"num1" binding:"required"`
	result chan int
}

func (task *SubTask) Process() {
	task.result <- task.Num1 - task.Num2
}

func (task *SubTask) Wait() int {
	return <-task.result
}
