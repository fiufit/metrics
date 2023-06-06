package handlers

import (
	"net/http"

	"github.com/fiufit/metrics/contracts"
	"github.com/fiufit/metrics/contracts/metrics"
	"github.com/fiufit/metrics/queue_services"
	"github.com/gin-gonic/gin"
)

type CreateMetric struct {
	metricsPublisher queue_services.MetricsPublisher
}

func NewCreateMetric(metricsPublisher queue_services.MetricsPublisher) CreateMetric {
	return CreateMetric{metricsPublisher: metricsPublisher}
}

func (h CreateMetric) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req metrics.CreateMetricRequest
		err := ctx.ShouldBindJSON(&req)
		validateErr := metrics.ValidateMetricTypes(req.MetricType, req.SubType)
		if err != nil || validateErr != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		h.metricsPublisher.PublishMetric(ctx, req)
	}
}
