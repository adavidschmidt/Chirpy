package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Kerfuffle this nonsense!",
			expected: "**** this nonsense!",
		},
		{
			input:    "This is a clean string",
			expected: "This is a clean string",
		},
		{
			input:    "Kerfuffle! This should stay",
			expected: "Kerfuffle! This should stay",
		},
		{
			input:    "Sharbert fornax kerfuffle",
			expected: "**** **** ****",
		},
	}

	for _, c := range cases {
		actual := helperCleanBody(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of '%s' is not as expected", actual)
		}
		if c.expected != actual {
			t.Errorf("Test failed '%s' does not equal '%s'", actual, c.expected)
		}
	}
}
