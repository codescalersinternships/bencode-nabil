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
			input:    "l4:spam4:eggse",
			expected: []interface{}{"spam", "eggs"},
			expectedError: false,
		},
		{
			name: "Dictionary with strings",
			input:    "d3:cow3:moo4:spam4:eggse",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"cow":"moo", "spam":"eggs",
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
			input:    "d3:cow3:moo4:spami123ee",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"cow":"moo", "spam":int64(123),
				},
			),
			expectedError: false,
		},
		{
			name: "Dictionary with strings and arrays",
			input:    "d3:cow3:moo4:spaml4:spam4:eggsee",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"cow":"moo", 
					"spam":[]interface {}{"spam", "eggs"},
				},
			),
			expectedError: false,
		},
		{
			name: "Dictionary with strings and dictionaries",
			input:    "d3:cow3:moo4:spamd3:cow3:moo4:spam4:eggsee",
			expected: map[interface {}]interface {}(
				map[interface {}]interface {}{
					"cow":"moo", 
					"spam":map[interface {}]interface {}{"cow":"moo", "spam":"eggs"},
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