package metrics

import (
	"time"

	"github.com/fiufit/metrics/models"
)

type GetMetricsRequest struct {
	MetricType string    `form:"type" binding:"required"`
	SubType    string    `form:"subtype"`
	From       time.Time `form:"from"`
	To         time.Time `form:"to"`
}

type GetMetricsResponse []models.Metric
