package monitoring

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	counter   *prometheus.CounterVec
	counterDB *prometheus.CounterVec
	latency   *prometheus.HistogramVec
	latencyDB *prometheus.HistogramVec
}

var (
	prom          *metrics
	dbName        = "db_requests_total"
	dbHelp        = "How many Database requests/calls processed, partitioned by status and table database"
	dbLatencyName = "db_requests_duration"
	dbLatencyHelp = "How long it took to process the database requests/calls, partitioned by status and table database"
	reqsName      = "http_requests_total"
	reqsHelp      = "How many HTTP requests processed, partitioned by status code, method and HTTP path."
	latencyName   = "http_request_duration_seconds"
	latencyHelp   = "How long it took to process the request, partitioned by status code, method and HTTP path."

	DefaultBuckets = []float64{0.3, 1.2, 5.0}
)

func NewPrometheus(name string) {
	str := []string{"code", "method", "path", "desc"}
	strDB := []string{"status", "table"}
	reqCounter := promauto.NewCounterVec(prometheus.CounterOpts{
		Help:        reqsHelp,
		Name:        reqsName,
		ConstLabels: prometheus.Labels{"service": name},
	}, str)

	reqLatency := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        latencyName,
		Help:        latencyHelp,
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     DefaultBuckets,
	}, str)

	reqCounterDB := promauto.NewCounterVec(prometheus.CounterOpts{
		Help:        dbHelp,
		Name:        dbName,
		ConstLabels: prometheus.Labels{"service": name, "database": "true"},
	}, strDB)

	reqLatencyDB := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        dbLatencyName,
		Help:        dbLatencyHelp,
		ConstLabels: prometheus.Labels{"service": name, "database": "true"},
		Buckets:     DefaultBuckets,
	}, strDB)

	prom = &metrics{
		counter: reqCounter,
		latency: reqLatency,
		// database
		counterDB: reqCounterDB,
		latencyDB: reqLatencyDB,
	}
}

func (r *metrics) Record(statusCode int, method, path, desc string, duration time.Duration) {
	code := strconv.Itoa(statusCode)

	r.counter.WithLabelValues(code, method, path, desc).Inc()
	r.latency.WithLabelValues(code, method, path, desc).Observe(float64(duration.Nanoseconds()) / 1000000000)
}

func (r *metrics) RecordDatabase(table string, start time.Time, err error) {
	status := "success"
	if err != nil {
		status = "failed"
	}

	since := time.Since(start)
	r.counterDB.WithLabelValues(status, table).Inc()
	r.latencyDB.WithLabelValues(status, table).Observe(float64(since.Nanoseconds()) / 1000000000)
}

func Prometheus() *metrics {
	return prom
}
