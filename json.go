package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processJSONKeyLowerFirst(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			newKey := strings.ToLower(string(key[0])) + key[1:]
			newMap[newKey] = processJSONKeyLowerFirst(value)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, value := range v {
			newSlice[i] = processJSONKeyLowerFirst(value)
		}
		return newSlice
	default:
		return v
	}
}

func processJSONKeyUpperFirstBytes(input []byte) ([]byte, error) {
	var data interface{}

	// Unmarshal the input JSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// Process the JSON
	processed := processJSONKeyLowerFirst(data)

	// Marshal the processed data back to JSON
	output, err := json.Marshal(processed)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	return output, nil
}

func removeEmptyValues(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			cleaned := removeEmptyValues(value)
			if cleaned != nil {
				newMap[key] = cleaned
			}
		}
		if len(newMap) > 0 {
			return newMap
		}
		return nil
	case []interface{}:
		newSlice := make([]interface{}, 0, len(v))
		for _, value := range v {
			cleaned := removeEmptyValues(value)
			if cleaned != nil {
				newSlice = append(newSlice, cleaned)
			}
		}
		if len(newSlice) > 0 {
			return newSlice
		}
		return nil
	case string:
		if v == "" {
			return nil
		}
		return v
	case nil:
		return nil
	default:
		return v
	}
}

func removeEmptyValuesBytes(input []byte) ([]byte, error) {
	var data interface{}
	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	cleaned := removeEmptyValues(data)

	result, err := json.Marshal(cleaned)
	if err != nil {
		return nil, fmt.Errorf("error marshaling cleaned JSON: %w", err)
	}

	return result, nil
}
