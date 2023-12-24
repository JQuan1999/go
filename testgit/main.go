package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cfn := context.WithCancel(context.Background())
	cal := NewCalculator(ctx, 10, 1024)
	go func() {
		cal.Run()
	}()
	go func() {
		ServerHttp(ctx, "127.0.0.1:8888", cal)
	}()
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	<-sign
	cfn()
	log.Println("exit")
}
