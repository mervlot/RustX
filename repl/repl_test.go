package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestReplExecution(t *testing.T) {
	input := "x := 10\nx + 5\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "RustX Starter v0.1") {
		t.Errorf("output does not contain banner: %s", output)
	}
	if !strings.Contains(output, "15") {
		t.Errorf("output does not contain '15': %s", output)
	}
}
