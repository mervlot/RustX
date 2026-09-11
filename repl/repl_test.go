package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestREPL(t *testing.T) {
	input := `x := 10
println(x)
exit
`
	in := strings.NewReader(input)
	var out bytes.Buffer

	Start(in, &out)

	output := out.String()
	if !strings.Contains(output, "10") {
		t.Fatalf("expected REPL output to contain '10', got=%q", output)
	}
}
