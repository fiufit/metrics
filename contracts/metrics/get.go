package metrics

import (
	"time"

	"github.com/fiufit/metrics/models"
)

type GetMetricsRequest struct {
	MetricType string `binding:"required"`
	SubType    string
	From       time.Time
	To         time.Time
}

type GetMetricsResponse []models.Metric
