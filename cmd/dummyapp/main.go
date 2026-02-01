// dummyapp is a fake web application that exposes Prometheus metrics
// simulating a real service with HTTP RED metrics, business KPIs,
// job queues, DB connections, and cache statistics.
package main

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_http_requests_total",
		Help: "Total HTTP requests.",
	}, []string{"method", "path", "status"})

	httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myapp_http_request_duration_seconds",
		Help:    "HTTP request latency.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	httpRequestsInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_http_requests_in_flight",
		Help: "Number of HTTP requests currently in flight.",
	})

	ordersCreatedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_orders_created_total",
		Help: "Total orders created.",
	}, []string{"status"})

	revenueTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_revenue_total",
		Help: "Total revenue.",
	}, []string{"currency"})

	usersActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_users_active",
		Help: "Number of currently active users.",
	})

	usersRegisteredTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_users_registered_total",
		Help: "Total user registrations.",
	})

	jobsProcessedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_jobs_processed_total",
		Help: "Total background jobs processed.",
	}, []string{"queue", "status"})

	jobsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myapp_jobs_duration_seconds",
		Help:    "Background job processing time.",
		Buckets: prometheus.DefBuckets,
	}, []string{"queue"})

	jobsQueueDepth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "myapp_jobs_queue_depth",
		Help: "Number of pending jobs in queue.",
	}, []string{"queue"})

	dbConnectionsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_db_connections_active",
		Help: "Number of active DB connections.",
	})

	dbConnectionsIdle = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_db_connections_idle",
		Help: "Number of idle DB connections.",
	})

	dbQueryDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myapp_db_query_duration_seconds",
		Help:    "DB query latency.",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation"})

	cacheHitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_cache_hits_total",
		Help: "Total cache hits.",
	})

	cacheMissesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_cache_misses_total",
		Help: "Total cache misses.",
	})

	errorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_errors_total",
		Help: "Total errors by type.",
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInFlight,
		ordersCreatedTotal,
		revenueTotal,
		usersActive,
		usersRegisteredTotal,
		jobsProcessedTotal,
		jobsDuration,
		jobsQueueDepth,
		dbConnectionsActive,
		dbConnectionsIdle,
		dbQueryDuration,
		cacheHitsTotal,
		cacheMissesTotal,
		errorsTotal,
	)
}

func main() {
	port := "3000"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Start background simulation goroutines.
	go simulateOrders()
	go simulateUsers()
	go simulateJobs()
	go simulateDBPool()
	go simulateCache()
	go simulateErrors()

	mux := http.NewServeMux()
	mux.HandleFunc("/", instrumentHandler("GET", "/", handleIndex))
	mux.HandleFunc("/api/users", instrumentHandler("GET", "/api/users", handleAPIUsers))
	mux.HandleFunc("/api/orders", instrumentHandler("POST", "/api/orders", handleAPIOrders))
	mux.HandleFunc("/api/search", instrumentHandler("GET", "/api/search", handleAPISearch))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "OK")
	})
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("dummyapp starting", "port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

// instrumentHandler wraps a handler to record HTTP metrics.
func instrumentHandler(method, path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		start := time.Now()
		next(w, r)
		duration := time.Since(start).Seconds()

		status := "200"
		// Simulate occasional errors.
		if rand.Float64() < 0.02 {
			status = "500"
		}

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}

func handleIndex(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Duration(5+rand.IntN(15)) * time.Millisecond)
	_, _ = fmt.Fprintln(w, "OK")
}

func handleAPIUsers(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Duration(10+rand.IntN(40)) * time.Millisecond)
	// Simulate DB query.
	dbQueryDuration.WithLabelValues("select").Observe(0.005 + rand.Float64()*0.02)
	_, _ = fmt.Fprintln(w, `{"users":[]}`)
}

func handleAPIOrders(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Duration(20+rand.IntN(80)) * time.Millisecond)
	// Simulate DB write.
	dbQueryDuration.WithLabelValues("insert").Observe(0.01 + rand.Float64()*0.03)
	_, _ = fmt.Fprintln(w, `{"order_id":"123"}`)
}

func handleAPISearch(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(time.Duration(30+rand.IntN(120)) * time.Millisecond)
	// Simulate DB query.
	dbQueryDuration.WithLabelValues("select").Observe(0.02 + rand.Float64()*0.05)
	_, _ = fmt.Fprintln(w, `{"results":[]}`)
}

// simulateOrders generates order and revenue metrics.
func simulateOrders() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// 70% completed, 20% failed, 10% cancelled
		r := rand.Float64()
		switch {
		case r < 0.7:
			ordersCreatedTotal.WithLabelValues("completed").Inc()
			revenueTotal.WithLabelValues("USD").Add(10 + rand.Float64()*90)
		case r < 0.9:
			ordersCreatedTotal.WithLabelValues("failed").Inc()
		default:
			ordersCreatedTotal.WithLabelValues("cancelled").Inc()
		}
	}
}

// simulateUsers generates user activity metrics.
func simulateUsers() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// Active users fluctuate with a sine wave pattern.
		t := float64(time.Now().Unix())
		active := 200 + 80*math.Sin(t/300) + float64(rand.IntN(30))
		usersActive.Set(active)

		// Occasional registrations.
		if rand.Float64() < 0.3 {
			usersRegisteredTotal.Inc()
		}
	}
}

// simulateJobs generates background job metrics.
func simulateJobs() {
	queues := []string{"email", "export", "cleanup"}
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for _, q := range queues {
			// Process 1-3 jobs per tick.
			n := 1 + rand.IntN(3)
			for range n {
				duration := 0.1 + rand.Float64()*2.0
				jobsDuration.WithLabelValues(q).Observe(duration)

				if rand.Float64() < 0.9 {
					jobsProcessedTotal.WithLabelValues(q, "success").Inc()
				} else {
					jobsProcessedTotal.WithLabelValues(q, "failure").Inc()
				}
			}
			// Queue depth fluctuates.
			depth := float64(rand.IntN(20))
			jobsQueueDepth.WithLabelValues(q).Set(depth)
		}
	}
}

// simulateDBPool generates DB connection pool metrics.
func simulateDBPool() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	maxConns := 20.0
	for range ticker.C {
		active := 5 + float64(rand.IntN(10))
		if active > maxConns {
			active = maxConns
		}
		idle := maxConns - active
		dbConnectionsActive.Set(active)
		dbConnectionsIdle.Set(idle)
	}
}

// simulateCache generates cache hit/miss metrics.
func simulateCache() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		hits := 5 + rand.IntN(15)
		misses := 1 + rand.IntN(5)
		cacheHitsTotal.Add(float64(hits))
		cacheMissesTotal.Add(float64(misses))
	}
}

// simulateErrors generates error metrics.
func simulateErrors() {
	types := []string{"timeout", "connection_refused", "internal", "validation"}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// Pick a random error type.
		t := types[rand.IntN(len(types))]
		errorsTotal.WithLabelValues(t).Inc()
	}
}
