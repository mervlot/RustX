package evaluator

import (
	"testing"

	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
)

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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"50 / 2 * 2 + 10", 60},
		{"10 % 3", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"true && true", true},
		{"true && false", false},
		{"false || true", true},
		{"!true", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestVarAndAssign(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"a := 10; a = a + 5; a;", 15},
		{"x := 0; for i := 0; i < 5; i++ { x = x + i; } x;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestRecursiveFunction(t *testing.T) {
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

func TestArrayAssignmentAndIndexing(t *testing.T) {
	input := `
arr := [10, 20, 30];
arr[1] = 99;
arr[1];
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 99)
}

func TestDivisionByZero(t *testing.T) {
	input := `10 / 0;`
	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected error object, got %T", evaluated)
	}
	if errObj.Message != "division by zero" {
		t.Errorf("wrong error message: %s", errObj.Message)
	}
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`len("hello")`, "5"},
		{`str_upper("rustx")`, "RUSTX"},
		{`str_lower("RUSTX")`, "rustx"},
		{`math_abs(-42)`, "42"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Inspect() != tt.expected {
			t.Errorf("expected=%s, got=%s", tt.expected, evaluated.Inspect())
		}
	}
}
