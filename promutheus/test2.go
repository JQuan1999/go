package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type myCollector struct {
	randomValue *prometheus.Desc
}

func (c *myCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.randomValue
}

// 实现collect方法 收集指标的值
func (c *myCollector) Collect(ch chan<- prometheus.Metric) {
	// 模拟收集指标的过程
	rand.Seed(time.Now().Unix())
	value := rand.Float64()
	ch <- prometheus.MustNewConstMetric(c.randomValue, prometheus.GaugeValue, value)
}

func newMyCollector() *myCollector {
	return &myCollector{
		randomValue: prometheus.NewDesc("random_value", "A random value between 0 and 1", nil, nil),
	}
}

var reg *prometheus.Registry = prometheus.NewRegistry()

func TestCollector() {
	collector := newMyCollector()
	reg.MustRegister(collector)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
