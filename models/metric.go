package models

import (
	"time"
)

type Metric struct {
	MetricType string    `json:"type" bson:"type"`
	SubType    string    `json:"subtype" bson:"subtype,omitempty"`
	DateTime   time.Time `json:"date_time" bson:"timestamp"`
}
