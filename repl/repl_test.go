package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestIsBalanced(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"x := 5", true},
		{"func main() {", false},
		{"func main() {\n    println(\"Hello\")\n}", true},
		{"arr := [1, 2, 3", false},
		{"arr := [1, 2, 3]", true},
		{"s := \"{ unclosed brace inside string \"", true},
	}

	for _, tt := range tests {
		result := isBalanced(tt.input)
		if result != tt.expected {
			t.Errorf("isBalanced(%q) got=%t, want=%t", tt.input, result, tt.expected)
		}
	}
}

func TestReplStart(t *testing.T) {
	input := "x := 10\nprintln(x)\nexit\n"
	in := strings.NewReader(input)
	var out bytes.Buffer

	Start(in, &out)

	output := out.String()
	if !strings.Contains(output, "RustX Starter v0.1") {
		t.Errorf("REPL banner missing")
	}
	if !strings.Contains(output, "10") {
		t.Errorf("REPL output did not contain '10'")
	}
}
