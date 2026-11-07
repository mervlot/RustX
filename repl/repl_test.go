package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplExecution(t *testing.T) {
	input := "x := 5\nx + 10\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "15") {
		t.Errorf("expected REPL output to contain 15, got %q", output)
	}
}

func TestReplMultiline(t *testing.T) {
	input := "func add(a int, b int) int {\nreturn a + b\n}\nadd(2, 3)\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "5") {
		t.Errorf("expected REPL output to contain 5, got %q", output)
	}
}
