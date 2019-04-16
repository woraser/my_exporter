---
title: Prometheus
date: 2019-04-01 16:40:00
tags: 系统运维
categories:
- 系统运维
- golang
- prometheus
---

## prometheus简述
    Prometheus是一个开源的完整监控解决方案
	[中文教程](https://yunlzheng.gitbook.io/prometheus-book/)
    [英文教程](https://prometheus.io/docs/introduction/overview/)
    
## 运行流程
    Server端通过exporter和pushgateway获取监控数据，存储到自带内置的时间序列数据库当中(TSDB)
    prometheus提供了内置ui和rest接口供外部调用，查看数据
 
## 自定义exporter
  exporter本质上是一个web服务，通过对外暴露路由来展示采集数据
通过prometheus提供的client_golang工具包，我们可以很方便的自定义一个exorter

    * step1
    
    自定义exporter struct
    ```
        type MyExporter struct {
        	namespace string
        	gaugeMetrics *prometheus.GaugeVec
        	duration prometheus.Gauge
        	totalScrapes prometheus.Counter
        }

    ```
        
    * step2
    
    实现promethes exporter的标准接口Describe和Collect
    ```
        func (e *MyExporter) Describe(ch chan<- *prometheus.Desc) {
        }
        func (e *MyExporter) Collect(ch chan<- *prometheus.Metric) {
        //collect logic
        }
    
    ```
    

    * step3
    
    将自定义exporter注册到prometheus中
    ```
    
    	registry := prometheus.NewRegistry()
    	registry.MustRegister(myExporter)
    	
    ```        
        
    * step4
    
    注册路由到http
    
    ```
        handle := promhttp.HandlerFor(registry,promhttp.HandlerOpts{})
    	http.Handle(*metricPath, handle)
    	log.Fatal(http.ListenAndServe(*listenAddress,nil))
    
    ```

    
        