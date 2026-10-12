package evaluator

import (
	"testing"

	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatAndStringExpression(t *testing.T) {
	strEval := testEval(`"hello " + "world"`)
	str, ok := strEval.(*object.String)
	if !ok || str.Value != "hello world" {
		t.Fatalf("string eval failed. got=%v", strEval)
	}

	fltEval := testEval(`3.5 + 2.5`)
	flt, ok := fltEval.(*object.Float)
	if !ok || flt.Value != 6.0 {
		t.Fatalf("float eval failed. got=%v", fltEval)
	}
}

func TestEvalControlFlowAndVariables(t *testing.T) {
	input := `
var x = 10;
y := 20;
var sum = 0;
for i := 0; i < 5; i++ {
	sum = sum + i;
}
sum;
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestEvalFunctionsAndRecursion(t *testing.T) {
	input := `
func factorial(n int) int {
	if n <= 1 {
		return 1;
	}
	return n * factorial(n - 1);
}
factorial(5);
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 120)
}

func TestEvalBuiltins(t *testing.T) {
	input := `
s := "RustX";
len(s);
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}
