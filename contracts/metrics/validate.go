package metrics

import "errors"

var validTypeSubTypes = map[string]struct{}{
	"register":         {},
	"login":            {},
	"blocked":          {},
	"password_recover": {},
	"location":         {},
	"new_training":     {},
}

func ValidateMetricTypes(metricType string, subType string) error {
	if _, ok := validTypeSubTypes[metricType]; !ok {
		return errors.New("Invalid metric type")
	}
	if (metricType == "register" || metricType == "login") && (subType != "mail" && subType != "federated_entity") {
		return errors.New("Invalid metric subtype")
	}

	return nil
}
