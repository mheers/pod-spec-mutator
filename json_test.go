package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple object",
			input:    `{"Name":"John","Age":30}`,
			expected: `{"name":"John","age":30}`,
		},
		{
			name:     "Nested object",
			input:    `{"Person":{"Name":"John","Age":30}}`,
			expected: `{"person":{"name":"John","age":30}}`,
		},
		{
			name:     "Array of objects",
			input:    `{"People":[{"Name":"John","Age":30},{"Name":"Jane","Age":25}]}`,
			expected: `{"people":[{"name":"John","age":30},{"name":"Jane","age":25}]}`,
		},
		{
			name:     "Mixed types",
			input:    `{"Name":"John","Age":30,"Hobbies":["Reading","Gaming"],"Address":{"City":"New York","ZipCode":"10001"}}`,
			expected: `{"name":"John","age":30,"hobbies":["Reading","Gaming"],"address":{"city":"New York","zipCode":"10001"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input interface{}
			err := json.Unmarshal([]byte(tt.input), &input)
			if err != nil {
				t.Fatalf("Failed to unmarshal input JSON: %v", err)
			}

			result := processJSONKeyLowerFirst(input)

			resultJSON, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}

			var expectedMap, resultMap map[string]interface{}
			err = json.Unmarshal([]byte(tt.expected), &expectedMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected JSON: %v", err)
			}
			err = json.Unmarshal(resultJSON, &resultMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal result JSON: %v", err)
			}

			if !reflect.DeepEqual(expectedMap, resultMap) {
				t.Errorf("processJSON() = %v, want %v", string(resultJSON), tt.expected)
			}
		})
	}
}

func TestProcessJSONBytes(t *testing.T) {
	result, err := processJSONKeyUpperFirstBytes([]byte(patchJSONTestUpperFirst))
	require.NoError(t, err)
	require.Equal(t, patchJSONTest, string(result))
}

func TestRemoveEmptyValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple object with empty value",
			input:    `{"a":"1","b":""}`,
			expected: `{"a":"1"}`,
		},
		{
			name:     "Nested object with empty values",
			input:    `{"a":{"b":"","c":"2"},"d":""}`,
			expected: `{"a":{"c":"2"}}`,
		},
		{
			name:     "Array with empty values",
			input:    `{"a":["1","","2"],"b":[]}`,
			expected: `{"a":["1","2"]}`,
		},
		{
			name:     "Complex nested structure",
			input:    `{"HostAliases":[{"Ip":"192.168.1.100","Hostnames":["foo.local"]}],"Containers":[{"Name":"postgres","ImagePullPolicy":"Never","Image":""}]}`,
			expected: `{"HostAliases":[{"Hostnames":["foo.local"],"Ip":"192.168.1.100"}],"Containers":[{"ImagePullPolicy":"Never","Name":"postgres"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var input interface{}
			err := json.Unmarshal([]byte(tt.input), &input)
			if err != nil {
				t.Fatalf("Failed to unmarshal input: %v", err)
			}

			result := removeEmptyValues(input)

			resultJSON, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("Failed to marshal result: %v", err)
			}

			var expectedMap, resultMap map[string]interface{}
			err = json.Unmarshal([]byte(tt.expected), &expectedMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected: %v", err)
			}
			err = json.Unmarshal(resultJSON, &resultMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal result: %v", err)
			}

			if !reflect.DeepEqual(expectedMap, resultMap) {
				t.Errorf("removeEmptyValues() = %v, want %v", string(resultJSON), tt.expected)
			}
		})
	}
}

func TestRemoveEmptyValuesBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple object with empty value",
			input:    `{"a":"1","b":""}`,
			expected: `{"a":"1"}`,
		},
		{
			name:     "Complex nested structure",
			input:    `{"HostAliases":[{"Ip":"192.168.1.100","Hostnames":["foo.local"]}],"Containers":[{"Name":"postgres","ImagePullPolicy":"Never","Image":""}]}`,
			expected: `{"HostAliases":[{"Hostnames":["foo.local"],"Ip":"192.168.1.100"}],"Containers":[{"ImagePullPolicy":"Never","Name":"postgres"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := removeEmptyValuesBytes([]byte(tt.input))
			if err != nil {
				t.Fatalf("removeEmptyValuesBytes() error = %v", err)
			}

			var expectedMap, resultMap map[string]interface{}
			err = json.Unmarshal([]byte(tt.expected), &expectedMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal expected: %v", err)
			}
			err = json.Unmarshal(result, &resultMap)
			if err != nil {
				t.Fatalf("Failed to unmarshal result: %v", err)
			}

			if !reflect.DeepEqual(expectedMap, resultMap) {
				t.Errorf("removeEmptyValuesBytes() = %v, want %v", string(result), tt.expected)
			}
		})
	}
}
