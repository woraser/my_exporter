# my_exporter
the demo for prometheus exporter


## create a exporter with prometheus client_golang

  * Step1:define exporter struct
  ```
      type MyExporter struct {
        namespace string
        gaugeMetrics *prometheus.GaugeVec
        duration prometheus.Gauge
        totalScrapes prometheus.Counter
      }

  ```

  * Step2: implement inteface:Describe, Collect
  ```
      func (e *MyExporter) Describe(ch chan<- *prometheus.Desc) {
      }
      func (e *MyExporter) Collect(ch chan<- *prometheus.Metric) {
      //collect logic
      }

  ```


  * Step3: register into prometheus
  ```

    registry := prometheus.NewRegistry()
    registry.MustRegister(myExporter)

  ```        

  * Step4: start with http

  ```
    handle := promhttp.HandlerFor(registry,promhttp.HandlerOpts{})
    http.Handle(*metricPath, handle)
    log.Fatal(http.ListenAndServe(*listenAddress,nil))

  ```
