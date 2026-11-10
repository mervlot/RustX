package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplExecution(t *testing.T) {
	input := "x := 5\nx + 10\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "15") {
		t.Fatalf("expected output to contain 15, got %q", output)
	}
}

func TestReplMultiline(t *testing.T) {
	input := "func square(n int) int {\nreturn n * n\n}\nsquare(4)\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "16") {
		t.Fatalf("expected output to contain 16, got %q", output)
	}
}
