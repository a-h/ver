package main

import "testing"
import "encoding/json"

func TestVersionNumbersCanBeAddedTogether(t *testing.T) {
	tests := []struct {
		start    Version
		delta    Version
		expected Version
	}{
		{
			start:    Version{1, 0, 0},
			delta:    Version{0, 0, 1},
			expected: Version{1, 0, 1},
		},
		{
			start:    Version{1, 2, 3},
			delta:    Version{1, 2, 3},
			expected: Version{2, 4, 6},
		},
	}

	for _, test := range tests {
		actual := test.start.Add(test.delta)

		if actual != test.expected {
			t.Errorf("for %v + %v; expected %v, but got %v", test.start, test.delta, test.expected, actual)
		}
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		input    Version
		expected string
	}{
		{
			input:    Version{1, 0, 0},
			expected: "1.0.0",
		},
		{
			input:    Version{1, 2, 3},
			expected: "1.2.3",
		},
	}

	for _, test := range tests {
		actual := test.input.String()

		if actual != test.expected {
			t.Errorf("for %v; expected '%v', but got %v", test.input, test.expected, actual)
		}
	}
}

func TestVersionJSON(t *testing.T) {
	tests := []struct {
		input    Version
		expected string
	}{
		{
			input:    Version{1, 0, 0},
			expected: "\"1.0.0\"",
		},
		{
			input:    Version{1, 2, 3},
			expected: "\"1.2.3\"",
		},
	}

	for _, test := range tests {
		b, err := json.Marshal(test.input)
		if err != nil {
			t.Errorf("Failed to marshal JSON: %v\n", err)
		}
		actual := string(b)

		if actual != test.expected {
			t.Errorf("for %v; expected '%v', but got %v", test.input, test.expected, actual)
		}
	}
}
