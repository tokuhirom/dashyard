package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP metrics.
var (
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dashyard_http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "dashyard_http_request_duration_seconds",
		Help:    "HTTP request latency in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

// Prometheus proxy metrics.
var (
	PrometheusQueryTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dashyard_prometheus_query_total",
		Help: "Total number of upstream Prometheus queries.",
	}, []string{"status"})

	PrometheusQueryDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "dashyard_prometheus_query_duration_seconds",
		Help:    "Upstream Prometheus query latency in seconds.",
		Buckets: prometheus.DefBuckets,
	})
)

// Dashboard metrics.
var (
	DashboardsLoaded = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dashyard_dashboards_loaded",
		Help: "Number of dashboard files currently loaded.",
	})

	DashboardReloadsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dashyard_dashboard_reloads_total",
		Help: "Total number of dashboard hot-reloads.",
	})
)
