package repl

import (
	"bytes"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
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
		t.Fatalf("expected version banner in output, got: %s", output)
	}
	if !strings.Contains(output, "15") {
		t.Fatalf("expected output '15', got: %s", output)
	}
}

func TestReplMultiline(t *testing.T) {
	input := "func test(a) {\nreturn a * 2\n}\ntest(21)\nexit\n"
	in := strings.NewReader(input)
	out := &bytes.Buffer{}

	Start(in, out)

	output := out.String()
	if !strings.Contains(output, "42") {
		t.Fatalf("expected '42', got: %s", output)
	}
}

func TestReplVariablesPreserved(t *testing.T) {
	env := object.NewEnvironment()

	l1 := lexer.New("a := 100")
	p1 := parser.New(l1)
	evaluator.Eval(p1.ParseProgram(), env)

	l2 := lexer.New("a + 25")
	p2 := parser.New(l2)
	res := evaluator.Eval(p2.ParseProgram(), env)

	intObj, ok := res.(*object.Integer)
	if !ok || intObj.Value != 125 {
		t.Fatalf("expected variable preserved with value 125, got: %v", res)
	}
}
