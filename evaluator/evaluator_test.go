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
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"1.5 + 2.5", 4.0},
		{"10.0 / 2.0", 5.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		floatObj, ok := evaluated.(*object.Float)
		if !ok {
			t.Fatalf("object is not Float. got=%T (%+v)", evaluated, evaluated)
		}
		if floatObj.Value != tt.expected {
			t.Errorf("object has wrong value. got=%g, want=%g", floatObj.Value, tt.expected)
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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true { 10 }", 10},
		{"if false { 10 }", nil},
		{"if 1 { 10 }", 10},
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 > 2 { 10 } else if 1 < 2 { 30 } else { 20 }", 30},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestVarAndAssignStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a int = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"a := 10; a;", 10},
		{"a := 10; a = 20; a;", 20},
		{"a := 5; a++; a;", 6},
		{"a := 5; a--; a;", 4},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObjectAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
add(5, 10)
`
	testIntegerObject(t, testEval(input), 15)
}

func TestRecursion(t *testing.T) {
	input := `
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n - 1)
}
factorial(5)
`
	testIntegerObject(t, testEval(input), 120)
}

func TestLoops(t *testing.T) {
	inputFor := `
sum := 0
for i := 0; i < 5; i++ {
	sum = sum + i
}
sum
`
	testIntegerObject(t, testEval(inputFor), 10)

	inputWhile := `
i := 0
sum := 0
while i < 5 {
	sum = sum + i
	i++
}
sum
`
	testIntegerObject(t, testEval(inputWhile), 10)
}

func TestStringAndArrays(t *testing.T) {
	inputStr := `
s := "hello " + "world"
s
`
	evaluated := testEval(inputStr)
	strObj, ok := evaluated.(*object.String)
	if !ok || strObj.Value != "hello world" {
		t.Fatalf("string test failed. got=%v", evaluated)
	}

	inputArr := `
arr := [10, 20, 30]
arr[1] = 99
arr[1]
`
	testIntegerObject(t, testEval(inputArr), 99)
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, 5},
		{`len([1, 2, 3])`, 3},
		{`to_upper("hello")`, "HELLO"},
		{`to_lower("WORLD")`, "world"},
		{`abs(-42)`, 42},
		{`pow(2, 3)`, 8},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			strObj, ok := evaluated.(*object.String)
			if !ok || strObj.Value != expected {
				t.Errorf("builtin string failed. got=%v, want=%s", evaluated, expected)
			}
		}
	}
}

func TestDivisionByZeroError(t *testing.T) {
	input := `10 / 0`
	evaluated := testEval(input)
	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Fatalf("expected error object, got=%T", evaluated)
	}
	if errObj.Message != "division by zero" {
		t.Errorf("expected 'division by zero', got=%q", errObj.Message)
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
