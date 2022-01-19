package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricService(t *testing.T) {
	metricService, err := NewPrometheusService()
	assert.NoError(t, err)
	assert.NotNil(t, metricService)
}
