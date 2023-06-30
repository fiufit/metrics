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

// Create Metric godoc
//
//	@Summary		Creates a new metric.
//	@Description	Creates a new metric. It must have a valid type + subtype.
//	@Tags			metrics
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string								true	"API Version"
//	@Param			payload								body		metrics.CreateMetricRequest	true	"Body params"
//	@Success		200									{object}	string	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Failure		500									{object}	contracts.ErrResponse
//	@Router			/{version}/metrics 	[post]
func (h CreateMetric) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req metrics.CreateMetricRequest
		err := ctx.ShouldBindJSON(&req)
		validateErr := metrics.ValidateMetricTypes(req.MetricType, &req.SubType)
		if err != nil || validateErr != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		if err := h.metricsPublisher.PublishMetric(ctx, req); err != nil {
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(""))
	}
}
