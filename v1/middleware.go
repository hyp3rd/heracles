package heracles

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

// Default buckets for latency metrics.
var defaultBucketsConfig = []float64{0.3, 1.2, 5.0}

// Constants for collector names.
const (
	RequestsCollectorName       = "chi_requests_total"
	LatencyCollectorName        = "chi_request_duration_seconds"
	DetailedErrorsCollectorName = "chi_detailed_errors_total"
	RequestSizeCollectorName    = "chi_request_size_bytes"
	ResponseSizeCollectorName   = "chi_response_size_bytes"
)

// Middleware is a handler that exposes prometheus metrics.
type Middleware struct {
	buckets             []float64
	requestsEnabled     bool
	latencyEnabled      bool
	requestSizeEnabled  bool
	responseSizeEnabled bool
	requests            *prometheus.CounterVec
	latency             *prometheus.HistogramVec
	requestSize         *prometheus.SummaryVec
	responseSize        *prometheus.SummaryVec
	customLabels        []string
	detailedErrorCount  *prometheus.CounterVec
}

// NewMiddleware creates a new instance of the Heracles Middleware with the specified name and options.
// The function initializes the Middleware struct with the provided options and sets default values for
// the buckets configuration if it is not provided. It also creates Prometheus metrics for tracking
// HTTP requests and latency if the corresponding options are enabled.
//
// Parameters:
// - name: The name of the service.
// - opts: Optional MiddlewareOption functions to customize the middleware behavior.
//
// Returns:
// - *Middleware: A pointer to the newly created Middleware instance.
// - error: An error if any occurred during the creation of the Middleware instance.
//
// Example usage:
//
//	middleware, err := NewMiddleware("my-service", WithRequestsEnabled(true), WithLatencyEnabled(false))
//	if err != nil {
//	    // Handle error
//	}
//	// Use the middleware instance
func NewMiddleware(name string, opts ...MiddlewareOption) (*Middleware, error) {
	m := &Middleware{}
	for _, opt := range opts {
		opt(m)
	}

	// if !m.latencyEnabled && !m.requestsEnabled {
	// 	return nil, errors.New("at least one of latency or requests must be enabled")
	// }

	if len(m.buckets) == 0 {
		m.buckets = defaultBucketsConfig
	}

	labelNames := append([]string{"code", "method", "path"}, m.customLabels...)

	if m.requestsEnabled {
		m.requests = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        RequestsCollectorName,
				Help:        "Number of HTTP requests partitioned by status code, method and HTTP path.",
				ConstLabels: prometheus.Labels{"service": name},
			}, labelNames)
		m.detailedErrorCount = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        DetailedErrorsCollectorName,
				Help:        "Detailed error counts partitioned by type, status code, method and HTTP path.",
				ConstLabels: prometheus.Labels{"service": name},
			}, []string{"type", "code", "method", "path"})
	}

	if m.latencyEnabled {
		m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        LatencyCollectorName,
			Help:        "Time spent on the request partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": name},
			Buckets:     m.buckets,
		}, labelNames)
	}

	if m.requestSizeEnabled {
		m.requestSize = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:        RequestSizeCollectorName,
			Help:        "Size of HTTP requests in bytes.",
			ConstLabels: prometheus.Labels{"service": name},
		}, labelNames)
	}

	if m.responseSizeEnabled {
		m.responseSize = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name:        ResponseSizeCollectorName,
			Help:        "Size of HTTP responses in bytes.",
			ConstLabels: prometheus.Labels{"service": name},
		}, labelNames)
	}

	return m, nil
}

// Handler returns a handler for the middleware pattern.
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		m.recordMetrics(ww, r, start)
	})
}

// recordMetrics records the metrics for the request.
func (m *Middleware) recordMetrics(ww middleware.WrapResponseWriter, r *http.Request, start time.Time) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return
	}

	rp := rctx.RoutePattern()
	duration := float64(time.Since(start).Seconds())
	status := strconv.Itoa(ww.Status())

	labels := m.collectLabels(r, status, rp)
	if m.requestsEnabled {
		m.requests.WithLabelValues(labels...).Inc()
	}

	if m.latencyEnabled {
		m.latency.WithLabelValues(labels...).Observe(duration)
	}

	if m.requestSizeEnabled {
		m.requestSize.WithLabelValues(labels...).Observe(float64(r.ContentLength))
	}

	if m.responseSizeEnabled {
		m.responseSize.WithLabelValues(labels...).Observe(float64(ww.BytesWritten()))
	}

	if m.requestsEnabled && ww.Status() >= 400 {
		errorType := "client_error"
		if ww.Status() >= 500 {
			errorType = "server_error"
		}
		m.detailedErrorCount.WithLabelValues(errorType, status, r.Method, rp).Inc()
	}
}

// collectLabels collects the labels for the metrics.
func (m *Middleware) collectLabels(r *http.Request, status, rp string) []string {
	labels := []string{status, r.Method, rp}
	for _, label := range m.customLabels {
		labels = append(labels, r.Header.Get(label)) // Custom labels from headers
	}
	return labels
}

// Collectors returns collectors for your own collector registry.
func (m *Middleware) Collectors() []prometheus.Collector {
	collectors := []prometheus.Collector{}
	if m.requestsEnabled {
		collectors = append(collectors, m.requests, m.detailedErrorCount)
	}
	if m.latencyEnabled {
		collectors = append(collectors, m.latency)
	}
	if m.requestSizeEnabled {
		collectors = append(collectors, m.requestSize)
	}
	if m.responseSizeEnabled {
		collectors = append(collectors, m.responseSize)
	}
	return collectors
}

// MustRegisterDefault registers collectors to DefaultRegisterer.
func (m *Middleware) MustRegisterDefault() {
	if len(m.Collectors()) == 0 {
		panic("collectors must be set")
	}
	prometheus.MustRegister(m.Collectors()...)
}
