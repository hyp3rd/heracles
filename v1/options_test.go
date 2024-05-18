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
	WithRequestsEnabled(true)(m)

	assert.True(t, m.requestsEnabled)
}

func TestWithLatencyEnabled(t *testing.T) {
	m := &Middleware{}
	WithLatencyEnabled(true)(m)

	assert.True(t, m.latencyEnabled)
}

func TestWithBuckets(t *testing.T) {
	m := &Middleware{}
	WithLatencyBuckets(0.2, 1.0, 4.0)(m)

	expectedBuckets := []float64{0.2, 1.0, 4.0}
	assert.Equal(t, expectedBuckets, m.buckets)
}
