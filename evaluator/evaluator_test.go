package evaluator

import (
	"testing"

	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
)

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
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"10 % 3", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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

func TestFunctionObject(t *testing.T) {
	input := `func factorial(n int) int {
	if n <= 1 {
		return 1;
	}
	return n * factorial(n - 1);
}
factorial(5);`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 120)
}

func TestForLoopAndVariables(t *testing.T) {
	input := `
var sum = 0;
for i := 1; i <= 5; i++ {
	sum = sum + i;
}
sum;
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, 5},
		{`len([1, 2, 3])`, 3},
		{`to_upper("rustx")`, "RUSTX"},
		{`to_lower("RUSTX")`, "rustx"},
		{`contains("hello world", "world")`, true},
		{`abs(-42)`, int64(42)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*object.String)
			if !ok || str.Value != expected {
				t.Errorf("expected string %s, got=%v", expected, evaluated)
			}
		case bool:
			b, ok := evaluated.(*object.Boolean)
			if !ok || b.Value != expected {
				t.Errorf("expected bool %t, got=%v", expected, evaluated)
			}
		}
	}
}

func TestDivisionByZeroError(t *testing.T) {
	input := `10 / 0;`
	evaluated := testEval(input)

	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("no error object returned. got=%T (%+v)", evaluated, evaluated)
	}

	if errObj.Message == "" {
		t.Errorf("error message is empty")
	}
}
