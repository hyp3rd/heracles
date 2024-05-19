package heracles

// MiddlewareOption is a functional option for the middleware.
type MiddlewareOption func(*Middleware)

// WithRequestsEnabled sets the requestsEnabled field of the middleware.
func WithRequestsEnabled() MiddlewareOption {
	return func(m *Middleware) {
		m.requestsEnabled = true
	}
}

// WithLatencyEnabled sets the latencyEnabled field of the middleware.
func WithLatencyEnabled() MiddlewareOption {
	return func(m *Middleware) {
		m.latencyEnabled = true
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

// WithRequestSizeEnabled sets the requestSizeEnabled field of the middleware.
func WithRequestSizeEnabled() MiddlewareOption {
	return func(m *Middleware) {
		m.requestSizeEnabled = true
	}
}

// WithResponseSizeEnabled sets the responseSizeEnabled field of the middleware.
func WithResponseSizeEnabled() MiddlewareOption {
	return func(m *Middleware) {
		m.responseSizeEnabled = true
	}
}
