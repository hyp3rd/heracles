package heracles

// MiddlewareOption is a functional option for the middleware.
type MiddlewareOption func(*Middleware)

// WithRequestsEnabled sets the requestsEnabled field of the middleware.
func WithRequestsEnabled(enabled bool) MiddlewareOption {
	return func(m *Middleware) {
		m.requestsEnabled = enabled
	}
}

// WithLatencyEnabled sets the latencyEnabled field of the middleware.
func WithLatencyEnabled(enabled bool) MiddlewareOption {
	return func(m *Middleware) {
		m.latencyEnabled = enabled
	}
}

// WithCustomLabels is a function that returns a MiddlewareOption which allows you to set custom labels for the Prometheus middleware.
// These custom labels can be used to add additional dimensions to the exported metrics.
func WithCustomLabels(customLabels ...string) MiddlewareOption {
	return func(m *Middleware) {
		m.customLabels = customLabels
	}
}

// WithLatencyBuckets sets the buckets for the latency histogram.
func WithLatencyBuckets(buckets ...float64) MiddlewareOption {
	return func(m *Middleware) {
		m.buckets = buckets
	}
}
