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
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
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

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"-2.5", -2.5},
		{"1.5 + 2.5", 4.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		result, ok := evaluated.(*object.Float)
		if !ok {
			t.Fatalf("object is not Float. got=%T (%+v)", evaluated, evaluated)
		}
		if result.Value != tt.expected {
			t.Errorf("object has wrong value. got=%g, want=%g", result.Value, tt.expected)
		}
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
		{"true == true", true},
		{"true != false", true},
		{"(1 < 2) && (2 < 3)", true},
		{"(1 > 2) || (2 < 3)", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
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

func TestFactorialRecursion(t *testing.T) {
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

func TestForLoopAndShorthand(t *testing.T) {
	input := `
var sum = 0;
for i := 0; i < 5; i++ {
	sum = sum + i;
}
sum;
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, int64(5)},
		{`to_upper("rustx")`, "RUSTX"},
		{`to_lower("RUSTX")`, "rustx"},
		{`contains("rustx starter", "starter")`, true},
		{`abs(-42)`, int64(42)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("expected String object, got=%T", evaluated)
			}
			if str.Value != expected {
				t.Errorf("wrong string value. got=%q, want=%q", str.Value, expected)
			}
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestDivisionByZeroError(t *testing.T) {
	input := `10 / 0`
	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("no error object returned. got=%T (%+v)", evaluated, evaluated)
	}
	if errObj.Err.Message != "division by zero" {
		t.Errorf("wrong error message. got=%q, want='division by zero'", errObj.Err.Message)
	}
}
