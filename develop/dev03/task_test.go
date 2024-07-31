package main

import (
	"os"
	"testing"
)

func TestSortLines(t *testing.T) {
	tests := []struct {
		input    []string
		column   int
		numeric  bool
		reverse  bool
		unique   bool
		month    bool
		humanNum bool
		expected []string
	}{
		{
			input:    []string{"a 2", "a 10", "a 1"},
			column:   1,
			numeric:  true,
			expected: []string{"a 1", "a 2", "a 10"},
		},
		{
			input:    []string{"a 1", "b 2", "a 1"},
			column:   0,
			unique:   true,
			expected: []string{"a 1", "b 2"},
		},

		{
			input:    []string{"a 1k", "a 2M", "a 500"},
			column:   1,
			humanNum: true,
			expected: []string{"a 500", "a 1k", "a 2M"},
		},
	}

	for _, test := range tests {
		result := sortLines(test.input, test.column, test.numeric, test.reverse, test.unique, test.month, test.humanNum)
		if len(result) != len(test.expected) {
			t.Errorf("Expected length %d, but got %d", len(test.expected), len(result))
		}
		for i := range result {
			if result[i] != test.expected[i] {
				t.Errorf("Expected %s, but got %s", test.expected[i], result[i])
			}
		}
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
