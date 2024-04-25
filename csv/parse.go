package csv

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var typeRank = map[string]int{
	"bool":    1,
	"int":     2,
	"float64": 3,
	"string":  4,
}

func minType(t1, t2 string) string {
	if typeRank[t1] > typeRank[t2] {
		return t1
	}
	return t2
}

func inferSliceTypeReduction(records []string) string {
	currentType := "bool" // Start with the most restrictive type

	for _, record := range records {
		value := record
		if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
			currentType = minType(currentType, "bool")
		} else if _, err := strconv.Atoi(value); err == nil {
			currentType = minType(currentType, "int")
		} else if _, err := strconv.ParseFloat(value, 64); err == nil {
			currentType = minType(currentType, "float64")
		} else {
			currentType = "string"
			break // once it's a string, it cannot go back to a more restrictive type
		}
	}

	return currentType
}

func stringToReflectType(t string) (reflect.Type, error) {
	switch t {
	case "int":
		return reflect.TypeOf([]int{}).Elem(), nil
	case "float64":
		return reflect.TypeOf([]float64{}).Elem(), nil
	case "bool":
		return reflect.TypeOf([]bool{}).Elem(), nil
	case "string":
		return reflect.TypeOf([]string{}).Elem(), nil
	default:
		return nil, fmt.Errorf("unsupported type '%s'", t)
	}
}

func convertStringToType(s string, typ reflect.Type) (reflect.Value, error) {
	switch typ.Kind() {
	case reflect.Int:
		if val, err := strconv.Atoi(s); err == nil {
			return reflect.ValueOf(val), nil
		}
	case reflect.Float64:
		if val, err := strconv.ParseFloat(s, 64); err == nil {
			return reflect.ValueOf(val), nil
		}
	case reflect.Bool:
		if val, err := strconv.ParseBool(s); err == nil {
			return reflect.ValueOf(val), nil
		}
	case reflect.String:
		return reflect.ValueOf(s), nil
	}
	return reflect.Value{}, fmt.Errorf("invalid type for conversion: %s", typ)
}
