package repl

import (
	"testing"
)

func TestIsBalanced(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"x := 5", true},
		{"func foo() {", false},
		{"func foo() {}", true},
		{"[1, 2, 3", false},
		{"[1, 2, 3]", true},
		{`"{"`, true},
		{`"{" {`, false},
		{`"{" {}`, true},
	}

	for _, tt := range tests {
		result := isBalanced(tt.input)
		if result != tt.expected {
			t.Errorf("isBalanced(%q) - expected=%t, got=%t", tt.input, tt.expected, result)
		}
	}
}
