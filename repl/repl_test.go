package repl

import "testing"

func TestIsBalanced(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"x := 5", true},
		{"func foo() {", false},
		{"func foo() {\n}", true},
		{"println(\"hello\"", false},
		{"println(\"hello\")", true},
		{"/* nested curly { */", true},
	}

	for _, tt := range tests {
		got := isBalanced(tt.input)
		if got != tt.expected {
			t.Errorf("isBalanced(%q) = %t; want %t", tt.input, got, tt.expected)
		}
	}
}
