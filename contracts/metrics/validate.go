package metrics

import "errors"

var validMetricTypes = map[string]struct{}{
	"register":         {},
	"login":            {},
	"blocked":          {},
	"password_recover": {},
	"location":         {},
	"new_training":     {},
	"training_tagged":  {},
}

func ValidateMetricTypes(metricType string, subType string) error {
	if _, ok := validMetricTypes[metricType]; !ok {
		return errors.New("Invalid metric type")
	}
	if (metricType == "register" || metricType == "login") && (subType != "mail" && subType != "federated_entity") {
		return errors.New("Invalid metric subtype")
	}

	return nil
}
