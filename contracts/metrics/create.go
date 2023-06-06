package metrics

type CreateMetricRequest struct {
	MetricType string `binding:"required"`
	SubType    string
}
