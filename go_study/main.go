package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 创建两个 Counter 指标
	counter1 := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "counter1",
		Help: "Sample counter 1",
	})
	counter2 := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "counter2",
		Help: "Sample counter 2",
	})

	// 注册两个 Counter 指标
	prometheus.MustRegister(counter1)
	prometheus.MustRegister(counter2)

	// 启动 HTTP 服务
	http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", nil)
}
