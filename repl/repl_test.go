package repl

import "testing"

func TestIsBalanced(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"x := 5", true},
		{"func main() {", false},
		{"func main() {\n}", true},
		{"var x = [1, 2, 3", false},
		{"var x = [1, 2, 3]", true},
		{"\"escaped { quote\"", true},
	}

	for _, tt := range tests {
		actual := isBalanced(tt.input)
		if actual != tt.expected {
			t.Errorf("isBalanced(%q) expected=%v, got=%v", tt.input, tt.expected, actual)
		}
	}
}
