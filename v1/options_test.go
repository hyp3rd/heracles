package heracles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithCustomLabels(t *testing.T) {
	m := &Middleware{}
	WithCustomLabels("label1", "label2")(m)

	expectedLabels := []string{"label1", "label2"}
	assert.Equal(t, expectedLabels, m.customLabels)
}

func TestWithRequestsEnabled(t *testing.T) {
	m := &Middleware{}
	WithRequestsEnabled()(m)

	assert.True(t, m.requestsEnabled)
}

func TestWithLatencyEnabled(t *testing.T) {
	m := &Middleware{}
	WithLatencyEnabled()(m)

	assert.True(t, m.latencyEnabled)
}

func TestWithBuckets(t *testing.T) {
	m := &Middleware{}
	WithLatencyBuckets(0.2, 1.0, 4.0)(m)

	expectedBuckets := []float64{0.2, 1.0, 4.0}
	assert.Equal(t, expectedBuckets, m.buckets)
}

func TestWithRequestSizeEnabled(t *testing.T) {
	m := &Middleware{}
	WithRequestSizeEnabled()(m)

	assert.True(t, m.requestSizeEnabled)
}

func TestWithResponseSizeEnabled(t *testing.T) {
	m := &Middleware{}
	WithResponseSizeEnabled()(m)

	assert.True(t, m.responseSizeEnabled)
}
