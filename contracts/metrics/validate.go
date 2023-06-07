package metrics

import "errors"

const anySubtype = "any"
const clearSubtype = ""

var validMetricTypeSubtypes = map[string][]string{
	"register":         {"mail", "federated_entity", ""},
	"login":            {"mail", "federated_entity", ""},
	"blocked":          {clearSubtype},
	"password_recover": {clearSubtype},
	"location":         {anySubtype},
	"new_training":     {clearSubtype},
	"training_tagged":  {"strength", "speed", "endurance", "lose weight", "gain weight", "sports", ""},
}

func ValidateMetricTypes(metricType string, subType *string) error {
	validSubtypes, typeExists := validMetricTypeSubtypes[metricType]
	if !typeExists {
		return errors.New("invalid metric type")
	}

	if validSubtypes[0] == anySubtype {
		return nil
	}

	if validSubtypes[0] == clearSubtype {
		*subType = clearSubtype
		return nil
	}

	for _, validSubType := range validSubtypes {
		if validSubType == *subType {
			return nil
		}
	}

	return errors.New("invalid metric subtype")
}
