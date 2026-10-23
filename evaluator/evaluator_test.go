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
		{"2 * (5 + 10)", 30},
		{"10 % 3", 1},
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

func TestVariablesAndControlFlow(t *testing.T) {
	input := `
var a = 5
b := 10
var sum = 0

for i := 1; i <= 5; i++ {
    sum = sum + i
}

if sum == 15 && a < b {
    sum
} else {
    0
}
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestFunctionAndRecursion(t *testing.T) {
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
s := "hello rustx"
u := to_upper(s)
arr := split(u, " ")
len(arr)
`
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 2)
}
