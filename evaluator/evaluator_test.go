package evaluator

import (
	"math"
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
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
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
	evaluated := testEval("3.14 + 2.0")
	fl, ok := evaluated.(*object.Float)
	if !ok {
		t.Fatalf("object is not Float. got=%T", evaluated)
	}
	if math.Abs(fl.Value-5.14) > 0.00001 {
		t.Fatalf("float value wrong. got=%g, want=5.14", fl.Value)
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
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"true && true", true},
		{"true && false", false},
		{"false || true", true},
		{"false || false", false},
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

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"func identity(x int) int { return x; }; identity(5);", 5},
		{"func double(x int) int { return x * 2; }; double(5);", 10},
		{"func add(x int, y int) int { return x + y; }; add(5, 5);", 10},
		{"func add(x int, y int) int { return x + y; }; add(5 + 5, add(5, 5));", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
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
	testIntegerObject(t, testEval(input), 120)
}

func TestForLoopAndVariables(t *testing.T) {
	input := `
sum := 0;
for i := 1; i <= 5; i++ {
	sum = sum + i;
}
sum;
`
	testIntegerObject(t, testEval(input), 15)
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, int64(5)},
		{`to_upper("hello")`, "HELLO"},
		{`to_lower("WORLD")`, "world"},
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
			if !ok || str.Value != expected {
				t.Fatalf("expected string %s, got %s", expected, evaluated.Inspect())
			}
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}
