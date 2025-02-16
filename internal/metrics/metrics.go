package metrics

import (
	"net/http"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TasksQueued = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "scheduler_tasks_queued",
		Help: "Number of tasks in the queue",
	})

	TasksProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "scheduler_tasks_processed_total",
		Help: "Total number of processed tasks",
	}, []string{"status"})

	RetryCount = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "scheduler_task_retries",
		Help:    "Number of retries per task",
		Buckets: []float64{0, 1, 2, 3, 5, 10},
	})
)

func Init(port int) {
	prometheus.MustRegister(TasksQueued, TasksProcessed, RetryCount)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()
}