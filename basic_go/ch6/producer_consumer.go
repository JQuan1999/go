package ch6

import (
	"context"
	"log"
	"math/rand"
	"time"
)

type Routine struct {
	ctx       context.Context
	nproducer int
	produers  []*Producer
	nconsumer int
	consumers []*Consumer

	taskQueue     chan [2]int
	taskQueueSize int
}

func NewRoutine(ctx context.Context, nproducer, nconsumer, taskQueueSize int) *Routine {
	var r Routine
	r.ctx = ctx
	r.taskQueueSize = taskQueueSize
	r.taskQueue = make(chan [2]int, taskQueueSize)

	r.nproducer = nproducer
	for i := 0; i < nproducer; i++ {
		p := NewProducer(ctx, r.taskQueue)
		go p.Run()
		r.produers = append(r.produers, p)
	}

	r.nconsumer = nconsumer
	for i := 0; i < nconsumer; i++ {
		c := NewConsumer(ctx, r.taskQueue)
		c.Run()
		r.consumers = append(r.consumers, c)
	}
	return &r
}

func (r *Routine) Run() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-r.ctx.Done():
			log.Println("routine closed...")
			return
		case <-ticker.C:
			log.Println("routine print every 5 seconds")
		}
	}
}

type Consumer struct {
	ctx   context.Context
	queue <-chan [2]int
}

func NewConsumer(ctx context.Context, queue <-chan [2]int) *Consumer {
	consumer := Consumer{ctx: ctx, queue: queue}
	return &consumer
}

func (c *Consumer) Run() {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("consumer closed...")
			return
		case nums := <-c.queue:
			result := nums[0] + nums[1]
			log.Println("result= ", result)
		}
	}
}

type Producer struct {
	ctx   context.Context
	queue chan<- [2]int
}

func NewProducer(ctx context.Context, queue chan<- [2]int) *Producer {
	return &Producer{ctx: ctx, queue: queue}
}

func (p *Producer) Run() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-p.ctx.Done():
			log.Println("producer closed...")
			return
		case <-ticker.C:
			r1 := rand.Intn(100)
			r2 := rand.Intn(100)
			p.queue <- [2]int{r1, r2}
		}
	}
}

func TestProduerConsumer() {
	ctx, cfn := context.WithTimeout(context.Background(), time.Minute*1)
	defer cfn()
	r := NewRoutine(ctx, 2, 3, 10)
	r.Run()
}
