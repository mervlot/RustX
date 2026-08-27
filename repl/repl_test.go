package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPL(t *testing.T) {
	input := `x := 10
y := 20
func add(a int, b int) int {
return a + b
}
add(x, y)
exit
`
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "30") {
		t.Fatalf("expected output to contain '30', got=%q", output)
	}
}
