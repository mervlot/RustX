package evaluator

import (
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"testing"
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
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"true && true", true},
		{"true && false", false},
		{"false || true", true},
		{"false || false", false},
		{"!true", false},
		{"!false", true},
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
		{"if 1 < 2 { 10 }", 10},
		{"if 1 > 2 { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			if evaluated != NULL {
				t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}

func TestForLoop(t *testing.T) {
	input := `
var sum = 0
for i := 1; i <= 5; i++ {
    sum = sum + i
}
sum
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestWhileLoop(t *testing.T) {
	input := `
var count = 5
var sum = 0
while count > 0 {
    sum = sum + count
    count--
}
sum
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestFunctionApplication(t *testing.T) {
	input := `
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

factorial(5)
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 120)
}

func TestStandardLibrary(t *testing.T) {
	input := `
s := "rustx"
len(s)
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)

	upperTest := `to_upper("rustx")`
	if res := testEval(upperTest); res.Inspect() != "RUSTX" {
		t.Fatalf("to_upper failed, got=%s", res.Inspect())
	}

	mathTest := `abs(-42)`
	testIntegerObject(t, testEval(mathTest), 42)
}

func TestArrayEvaluation(t *testing.T) {
	input := `
arr := [10, 20, 30]
arr[1] = 99
arr[1]
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 99)
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
