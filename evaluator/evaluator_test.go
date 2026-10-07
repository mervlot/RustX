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
	flt, ok := evaluated.(*object.Float)
	if !ok {
		t.Fatalf("object is not Float. got=%T (%+v)", evaluated, evaluated)
	}
	if math.Abs(flt.Value-5.14) > 0.00001 {
		t.Errorf("float value wrong. got=%g, want=5.14", flt.Value)
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
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"true && true", true},
		{"true && false", false},
		{"false || true", true},
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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", int64(10)},
		{"if false { 10 }", nil},
		{"if 1 < 2 { 10 }", int64(10)},
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", int64(20)},
		{"if 1 < 2 { 10 } else { 20 }", int64(10)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int64)
		if ok {
			testIntegerObject(t, evaluated, integer)
		} else {
			if evaluated != NULL {
				t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}

func TestFunctionApplicationAndRecursion(t *testing.T) {
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

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, int64(5)},
		{`to_upper("rustx")`, "RUSTX"},
		{`to_lower("RUSTX")`, "rustx"},
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
				t.Errorf("object is not String. got=%T", evaluated)
			} else if str.Value != expected {
				t.Errorf("string wrong. got=%s, want=%s", str.Value, expected)
			}
		}
	}
}

func TestDivisionByZeroError(t *testing.T) {
	evaluated := testEval("10 / 0")
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected Error object, got=%T", evaluated)
	}
	if errObj.Message != "division by zero" {
		t.Errorf("error message wrong. got=%q", errObj.Message)
	}
}
