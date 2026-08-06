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

func TestIntegerEvaluation(t *testing.T) {
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

func TestFloatEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10.5", 10.5},
		{"-1.5", -1.5},
		{"1.0 + 2.5", 3.5},
		{"5.5 / 2.0", 2.75},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestBooleanEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestShortCircuitLogical(t *testing.T) {
	// If false && anything, the second arg is not evaluated. If it had error, it doesn't run.
	// Since x is not defined, calling x should cause "identifier not found" error.
	// But false && x returns false immediately without error!
	tests := []struct {
		input    string
		expected bool
	}{
		{"false && x", false},
		{"true || x", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		input       string
		expectedErr string
	}{
		{"5 / 0", "division by zero"},
		{"5.5 / 0.0", "division by zero"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for %q. got=%T", tt.input, evaluated)
			continue
		}
		if errObj.Message != tt.expectedErr {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedErr, errObj.Message)
		}
	}
}

func TestVarAndShorthandAssignments(t *testing.T) {
	input := `
var x = 5;
y := 10;
x = x + y;
x;
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestRecursiveFunctions(t *testing.T) {
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

func TestForLoop(t *testing.T) {
	input := `
sum := 0;
for i := 1; i <= 5; i++ {
	sum = sum + i;
}
sum;
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestArrayIndexing(t *testing.T) {
	input := `
arr := [1, 2, 3];
arr[1];
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 2)
}

func TestStringIndexing(t *testing.T) {
	input := `
str := "hello";
str[1];
`
	evaluated := testEval(input)
	strObj, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("not a String object. got=%T", evaluated)
	}
	if strObj.Value != "e" {
		t.Errorf("expected 'e', got %q", strObj.Value)
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

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%g, want=%g", result.Value, expected)
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
