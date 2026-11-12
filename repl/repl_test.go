package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplExecution(t *testing.T) {
	input := "x := 10\nx * 2\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "20") {
		t.Fatalf("expected output to contain 20, got:\n%s", output)
	}
}

func TestReplMultiline(t *testing.T) {
	input := "func square(n int) int {\nreturn n * n\n}\nsquare(4)\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "16") {
		t.Fatalf("expected output to contain 16, got:\n%s", output)
	}
}
