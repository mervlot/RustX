package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPL(t *testing.T) {
	input := `x := 10
println(x)
func add(a int, b int) int {
    return a + b
}
add(x, 5)
exit
`
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()

	if !strings.Contains(output, "RustX Starter v0.1 REPL") {
		t.Errorf("REPL banner missing. Output:\n%s", output)
	}

	if !strings.Contains(output, "10") {
		t.Errorf("Expected '10' in output. Output:\n%s", output)
	}

	if !strings.Contains(output, "15") {
		t.Errorf("Expected '15' in output. Output:\n%s", output)
	}
}
