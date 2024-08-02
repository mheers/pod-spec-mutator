package main

import (
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func comparePods(original, patch corev1.Pod) ([]map[string]interface{}, error) {
	var ops []map[string]interface{}

	originalValue := reflect.ValueOf(original)
	patchValue := reflect.ValueOf(patch)

	for i := 0; i < patchValue.NumField(); i++ {
		field := patchValue.Type().Field(i)
		originalFieldValue := originalValue.FieldByName(field.Name)
		patchFieldValue := patchValue.Field(i)

		if !patchFieldValue.IsZero() && !reflect.DeepEqual(originalFieldValue.Interface(), patchFieldValue.Interface()) {
			path := fmt.Sprintf("/%s", getJSONFieldName(field))
			op := map[string]interface{}{
				"op":    "replace",
				"path":  path,
				"value": patchFieldValue.Interface(),
			}
			ops = append(ops, op)
		}
	}

	return ops, nil
}

func getJSONFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}
	return strings.Split(jsonTag, ",")[0]
}
