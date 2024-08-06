package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func processJSON(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for key, value := range v {
			newKey := strings.ToLower(string(key[0])) + key[1:]
			newMap[newKey] = processJSON(value)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for i, value := range v {
			newSlice[i] = processJSON(value)
		}
		return newSlice
	default:
		return v
	}
}

func processJSONBytes(input []byte) ([]byte, error) {
	var data interface{}

	// Unmarshal the input JSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// Process the JSON
	processed := processJSON(data)

	// Marshal the processed data back to JSON
	output, err := json.Marshal(processed)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	return output, nil
}
