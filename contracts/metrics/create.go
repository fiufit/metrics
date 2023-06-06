package metrics

type CreateMetricRequest struct {
	MetricType string `binding:"required" json:"type"`
	SubType    string `json:"subtype"`
}
