package prometheus

import (
	"os"
	"os/signal"
	"sync"
)

import (
	"github.com/lethexixin/go-funcs/common/graceful"
)

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricData struct {
	Flag     string
	Function string
}

type MetricHisData struct {
	Flag     string
	Function string
	Duration int
}

var (
	wg sync.WaitGroup
)

func NewCounter(name, help string, labelName []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labelName,
	)
}

func NewHistogram(name, help string, labelName []string, buckets []float64) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    help,
			Buckets: buckets,
		},
		labelName,
	)
}

func CollectMetricsHistogram(dataHisChan chan MetricHisData, metric *prometheus.HistogramVec) {
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, graceful.ShutdownSignals...)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case data, ok := <-dataHisChan:
					if ok {
						metric.With(prometheus.Labels{"flag": data.Flag, "function": data.Function}).Observe(float64(data.Duration))
					}
				case <-signChan:
					for d := range dataHisChan {
						metric.With(prometheus.Labels{"flag": d.Flag, "function": d.Function}).Observe(float64(d.Duration))
					}
					break
				}
			}
		}()
	}
	wg.Wait()
}

func CollectMetrics(dataChan chan MetricData, metric *prometheus.CounterVec) {
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, graceful.ShutdownSignals...)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case data, ok := <-dataChan:
					if ok {
						metric.With(prometheus.Labels{"flag": data.Flag, "function": data.Function}).Inc()
					}
				case <-signChan:
					for d := range dataChan {
						metric.With(prometheus.Labels{"flag": d.Flag, "function": d.Function}).Inc()
					}
					break
				}
			}
		}()
	}
	wg.Wait()
}
