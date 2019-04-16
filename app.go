package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
	"time"
)

// define  exporter
// 4 type for metric:Counter,Gauge,Histogram,Summary
type MyExporter struct {
	namespace string
	gaugeMetrics *prometheus.GaugeVec
	duration prometheus.Gauge
	totalScrapes prometheus.Counter
}

// implement interface Describe(ch chan<- *prometheus.Desc)
func (e *MyExporter) Describe(ch chan<- *prometheus.Desc) {
	//创建指定监控项的描述
	e.gaugeMetrics.Describe(ch)
	ch <- e.totalScrapes.Desc()
}

// implement interface Collect(ch chan<- prometheus.Metric)
// the function where handle logic for collect monitor data
// step1: define metric and collect labels and data
// step2: put metric to prometheus.Metric
func (e *MyExporter) Collect(ch chan<- prometheus.Metric) {
	go func(exp *MyExporter) {
		//set val
		exp.totalScrapes.Inc()
		exp.gaugeMetrics.WithLabelValues(strconv.Itoa(time.Now().Nanosecond()), "up").Set(1)
		exp.duration.Set(float64(time.Since(startTime).Nanoseconds()))
	}(e)

	//collect data: metric.Collect(ch) or ch <- metric
	e.duration.Collect(ch)
	e.gaugeMetrics.Collect(ch)
	e.totalScrapes.Collect(ch)
}

var (
	nameSpace = flag.String("name_space","private","the name space for metric")
	listenAddress     = flag.String("web.listen-address", ":9121", "Address to listen on for web interface and telemetry.")
	metricPath        = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	myExporter = &MyExporter{}
	startTime time.Time
)

func init() {
	flag.Parse()
	startTime = time.Now()

	myExporter = &MyExporter{
		totalScrapes:prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: *nameSpace,
			Name: "scrapes_total_num",
			Help: "the total num for scrapes"}),
		gaugeMetrics:prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: *nameSpace,
			Name:"gauge_test",
			Help:"the test for gauge",
		},[]string{"time","status"}),
		duration:prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: *nameSpace,
			Name:"duration_time",
			Help:"the nanoseconds time for duration",
		}),
	}
}


func main() {

	registry := prometheus.NewRegistry()
	registry.MustRegister(myExporter)
	handle := promhttp.HandlerFor(registry,promhttp.HandlerOpts{})
	http.Handle(*metricPath, handle)
	log.Fatal(http.ListenAndServe(*listenAddress,nil))
}