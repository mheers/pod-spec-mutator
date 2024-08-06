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

			result := processJSON(input)

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
	result, err := processJSONBytes([]byte(patchJSONTestUpperFirst))
	require.NoError(t, err)
	require.Equal(t, patchJSONTest, string(result))
}
