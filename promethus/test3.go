package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProxyInst struct {
	address string
	port    string
}

var key string = "dbproxy_value"

func (proxy *ProxyInst) labels() []string {
	return []string{"address", "port"}
}

func (proxy *ProxyInst) labelValues() []string {
	return []string{proxy.address, proxy.port}
}

type myProxyCollector struct {
	proxy ProxyInst
	desc  *prometheus.Desc
}

func NewProxyCollector(address string, port string) *myProxyCollector {
	var collector myProxyCollector
	collector.proxy.address = address
	collector.proxy.port = port
	return &collector
}

// 自定义collector实现了collector接口, 创建含有address和port标签的Desc
func (c *myProxyCollector) Describe(ch chan<- *prometheus.Desc) {
	c.desc = prometheus.NewDesc(prometheus.BuildFQName("test", "", key), key, c.proxy.labels(), nil)
	ch <- c.desc
}

// 写入与desc对应的metric
func (c *myProxyCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, 0.1, c.proxy.labelValues()...)
}

func TestCreateMetric() {
	collector := NewProxyCollector("1.2.3.4", "1234")
	prometheus.MustRegister(collector)
	http.Handle("/metric", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", nil)
}

// 测试
// curl -X GET 'http://127.0.0.1:8080/metrics'
// # HELP test_dbproxy_value dbproxy_value
// # TYPE test_dbproxy_value gauge
// test_dbproxy_value{address="1.2.3.4",port="1234"} 0.1
