package handlers

import (
	"net/http"

	"github.com/fiufit/metrics/contracts"
	"github.com/fiufit/metrics/contracts/metrics"
	"github.com/fiufit/metrics/repositories"
	"github.com/gin-gonic/gin"
)

type GetMetrics struct {
	metrics repositories.Metrics
}

func NewGetMetrics(metrics repositories.Metrics) GetMetrics {
	return GetMetrics{metrics: metrics}
}

// Get Metrics godoc
//
//	@Summary		Gets metrics.
//	@Description	Gets metrics with type/subtype/date filters.
//	@Tags			metrics
//	@Accept			json
//	@Produce		json
//	@Param			version								path		string								true	"API Version"
//	@Param			type								query		string								true	"Metrics type"
//	@Param			subtype								query		string								false	"Metrics subtype"
//	@Param			from								query		string								false	"Starting date"
//	@Param			to									query		string								false	"Ending date"
//	@Success		200									{object}	metrics.GetMetricsResponse	"Important Note: OK responses are wrapped in {"data": ... }"
//	@Failure		400									{object}	contracts.ErrResponse
//	@Failure		500									{object}	contracts.ErrResponse
//	@Router			/{version}/metrics 	[get]
func (h GetMetrics) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req metrics.GetMetricsRequest
		err := ctx.ShouldBindQuery(&req)
		validateErr := metrics.ValidateMetricTypes(req.MetricType, &req.SubType)
		if err != nil || validateErr != nil {
			ctx.JSON(http.StatusBadRequest, contracts.FormatErrResponse(contracts.ErrBadRequest))
			return
		}

		res, err := h.metrics.Get(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, contracts.FormatErrResponse(contracts.ErrInternal))
			return
		}

		ctx.JSON(http.StatusOK, contracts.FormatOkResponse(res))
	}
}
