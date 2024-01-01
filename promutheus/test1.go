package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	cpuTemp    prometheus.Gauge
	hdFailures *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp:    prometheus.NewGauge(prometheus.GaugeOpts{Name: "cpu_temperature_celsius", Help: "Current temperature of the CPU."}),
		hdFailures: prometheus.NewCounterVec(prometheus.CounterOpts{Name: "hd_errors_total", Help: "number of hard-disk errors."}, []string{"device"}),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.hdFailures)
	return m
}

func Test() {
	// 创建一个非全局的registry
	reg := prometheus.NewRegistry()

	// 创建新的metric并用定制的reg注册指标
	m := NewMetrics(reg)
	m.cpuTemp.Set(65.3)
	m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// 暴露指标
	// promhttp.Handler返回一个默认的http.Handler for prometheus.DefaultGatherer 封装过的函数
	http.Handle("/metrics1", promhttp.Handler())
	// HandlerFor返回一个未封装过的http.Handler
	http.Handle("/metrics2", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
