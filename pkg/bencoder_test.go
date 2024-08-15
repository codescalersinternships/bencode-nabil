package bencoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestDecoder(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		expectedError bool
	}{
		{
			name:     "Simple string",
			input:    "4:spam",
			expected: "spam",
			expectedError: false,
		},
		{
			name:     "Integer",
			input:    "i123e",
			expected: int64(123),
			expectedError: false,
		},
		{
			name:     "List of strings",
			input:    "l4:spam4:foooe",
			expected: []interface{}{"spam", "fooo"},
			expectedError: false,
		},
		{
			name: "Dictionary with strings",
			input:    "d3:bar3:moo4:spam4:foooe",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"bar":"moo", "spam":"fooo",
				},
			),
			expectedError: false,
		},
		{
			name:     "Invalid bencoded string",
			input:    "invalid",
			expected: nil,
			expectedError: true,
		},
		{
			name: "Dictionary with strings and integers",
			input:    "d3:bar3:moo4:spami123ee",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"bar":"moo", "spam":int64(123),
				},
			),
			expectedError: false,
		},
		{
			name: "Dictionary with strings and arrays",
			input:    "d3:bar3:moo4:spaml4:spam4:foooee",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"bar":"moo", 
					"spam":[]interface {}{"spam", "fooo"},
				},
			),
			expectedError: false,
		},
		{
			name: "Dictionary with strings and dictionaries",
			input:    "d3:bar3:moo4:spamd3:bar3:moo4:spam4:foooee",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"bar":"moo", 
					"spam":map[interface {}]interface {}{"bar":"moo", "spam":"fooo"},
				},
			),
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Decoder(test.input)
			if test.expectedError {
				assert.Error(t, err)
			}
			assert.Equal(t, test.expected, got)
		})
	}
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expected      string
		expectedError bool
	}{
		{
			name:     "Simple string",
			input:    "spam",
			expected: "4:spam",
			expectedError: false,
		},
		{
			name:     "Integer",
			input:    int64(123),
			expected: "i123e",
			expectedError: false,
		},
		{
			name:     "List of strings",
			input:    []interface{}{"spam", "fooo"},
			expected: "l4:spam4:foooe",
			expectedError: false,
		},
		{
			name: "Dictionary with strings",
			input: map[interface{}]interface{}{
				"bar": "moo", "spam": "fooo",
			},
			expected:      "d3:bar3:moo4:spam4:foooe",
			expectedError: false,
		},
		{
			name:          "Unsupported type",
			input:         3.14,
			expected:      "",
			expectedError: true,
		},
		{
			name: "Dictionary with strings and integers",
			input: map[interface{}]interface{}{
				"bar": "moo", "spam": int64(123),
			},
			expected:      "d3:bar3:moo4:spami123ee",
			expectedError: false,
		},
		{
			name: "Dictionary with strings and arrays",
			input: map[interface{}]interface{}{
				"bar": "moo", 
				"spam": []interface{}{"spam", "fooo"},
			},
			expected:      "d3:bar3:moo4:spaml4:spam4:foooee",
			expectedError: false,
		},
		{
			name: "Dictionary with strings and dictionaries",
			input: map[interface{}]interface{}{
				"bar": "moo",
				"spam": map[interface{}]interface{}{
					"bar": "moo", "spam": "fooo",
				},
			},
			expected:      "d3:bar3:moo4:spamd3:bar3:moo4:spam4:foooee",
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Encoder(test.input)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, string(got))
			}
		})
	}
}