package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
	}

	for _, test := range tests {
		result, err := unpackString(test.input)
		if test.err {
			if err == nil {
				t.Errorf("Expected error for input %s, but got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("For input %s, expected %s, but got %s", test.input, test.expected, result)
			}
		}
	}
}
