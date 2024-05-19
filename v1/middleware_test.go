package heracles

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

// TestNewMiddleware tests the creation of middleware with different options.
func TestNewMiddleware(t *testing.T) {
	m, err := NewMiddleware("test_service", WithRequestsEnabled(false), WithLatencyEnabled(false))
	assert.NoError(t, err)
	assert.Nil(t, m.requests)
	assert.Nil(t, m.latency)

	m, err = NewMiddleware("test_service", WithRequestsEnabled(true), WithLatencyEnabled(true))
	assert.NoError(t, err)
	assert.NotNil(t, m.requests)
	assert.NotNil(t, m.latency)
}

// TestMiddlewareHandler tests the middleware handler to ensure metrics are recorded.
func TestMiddlewareHandler(t *testing.T) {
	reg := prometheus.NewRegistry()
	m, err := NewMiddleware("test_service", WithRequestsEnabled(true), WithLatencyEnabled(true), WithRequestSizeEnabled(true), WithResponseSizeEnabled(true))
	assert.NoError(t, err)

	reg.MustRegister(m.Collectors()...)

	r := chi.NewRouter()
	r.Use(m.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", strings.NewReader("test request body"))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	expectedRequests := `
		# HELP chi_requests_total Number of HTTP requests partitioned by status code, method and HTTP path.
		# TYPE chi_requests_total counter
		chi_requests_total{code="200",method="GET",path="/",service="test_service"} 1
	`
	expectedLatency := `
		# HELP chi_request_duration_seconds Time spent on the request partitioned by status code, method and HTTP path.
		# TYPE chi_request_duration_seconds histogram
		chi_request_duration_seconds_bucket{code="200",method="GET",path="/",service="test_service",le="0.3"} 1
		chi_request_duration_seconds_bucket{code="200",method="GET",path="/",service="test_service",le="1.2"} 1
		chi_request_duration_seconds_bucket{code="200",method="GET",path="/",service="test_service",le="5"} 1
		chi_request_duration_seconds_bucket{code="200",method="GET",path="/",service="test_service",le="+Inf"} 1
		chi_request_duration_seconds_count{code="200",method="GET",path="/",service="test_service"} 1
	`
	expectedRequestSize := `
		# HELP chi_request_size_bytes Size of HTTP requests in bytes.
		# TYPE chi_request_size_bytes summary
		chi_request_size_bytes_sum{code="200",method="GET",path="/",service="test_service"} 17
		chi_request_size_bytes_count{code="200",method="GET",path="/",service="test_service"} 1
	`
	expectedResponseSize := `
		# HELP chi_response_size_bytes Size of HTTP responses in bytes.
		# TYPE chi_response_size_bytes summary
		chi_response_size_bytes_sum{code="200",method="GET",path="/",service="test_service"} 0
		chi_response_size_bytes_count{code="200",method="GET",path="/",service="test_service"} 1
	`

	err = testutil.CollectAndCompare(m.requests, strings.NewReader(expectedRequests))
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(m.latency, strings.NewReader(expectedLatency), "chi_request_duration_seconds_bucket", "chi_request_duration_seconds_count")
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(m.requestSize, strings.NewReader(expectedRequestSize))
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(m.responseSize, strings.NewReader(expectedResponseSize))
	assert.NoError(t, err)
}

// TestDetailedErrorMetrics tests the recording of detailed error metrics.
func TestDetailedErrorMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()
	m, err := NewMiddleware("test_service", WithRequestsEnabled(true))
	assert.NoError(t, err)

	reg.MustRegister(m.Collectors()...)

	r := chi.NewRouter()
	r.Use(m.Handler)
	r.Get("/client_error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	r.Get("/server_error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	req := httptest.NewRequest("GET", "/client_error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	req = httptest.NewRequest("GET", "/server_error", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedDetailedErrors := `
		# HELP chi_detailed_errors_total Detailed error counts partitioned by type, status code, method and HTTP path.
		# TYPE chi_detailed_errors_total counter
		chi_detailed_errors_total{code="400",method="GET",path="/client_error",service="test_service",type="client_error"} 1
		chi_detailed_errors_total{code="500",method="GET",path="/server_error",service="test_service",type="server_error"} 1
	`
	err = testutil.CollectAndCompare(m.detailedErrorCount, strings.NewReader(expectedDetailedErrors))
	assert.NoError(t, err)
}

// TestCustomLabels tests the middleware with custom labels.
func TestCustomLabels(t *testing.T) {
	reg := prometheus.NewRegistry()
	m, err := NewMiddleware("test_service", WithRequestsEnabled(true), WithLatencyEnabled(true), WithCustomLabels("X_Custom_Header"))
	assert.NoError(t, err)

	reg.MustRegister(m.Collectors()...)

	r := chi.NewRouter()
	r.Use(m.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X_Custom_Header", "test-value")
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X_Custom_Header", "test-value")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	expectedRequests := `
		# HELP chi_requests_total Number of HTTP requests partitioned by status code, method and HTTP path.
		# TYPE chi_requests_total counter
		chi_requests_total{X_Custom_Header="test-value",code="200",method="GET",path="/",service="test_service"} 1
	`
	expectedLatency := `
		# HELP chi_request_duration_seconds Time spent on the request partitioned by status code, method and HTTP path.
		# TYPE chi_request_duration_seconds histogram
		chi_request_duration_seconds_bucket{X_Custom_Header="test-value",code="200",method="GET",path="/",service="test_service",le="0.3"} 1
		chi_request_duration_seconds_bucket{X_Custom_Header="test-value",code="200",method="GET",path="/",service="test_service",le="1.2"} 1
		chi_request_duration_seconds_bucket{X_Custom_Header="test-value",code="200",method="GET",path="/",service="test_service",le="5"} 1
		chi_request_duration_seconds_bucket{X_Custom_Header="test-value",code="200",method="GET",path="/",service="test_service",le="+Inf"} 1
		chi_request_duration_seconds_count{X_Custom_Header="test-value",code="200",method="GET",path="/",service="test_service"} 1
	`

	err = testutil.CollectAndCompare(m.requests, strings.NewReader(expectedRequests))
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(m.latency, strings.NewReader(expectedLatency), "chi_request_duration_seconds_bucket", "chi_request_duration_seconds_count")
	assert.NoError(t, err)
}
