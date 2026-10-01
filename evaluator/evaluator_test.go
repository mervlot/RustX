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
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
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

func TestFunctionAndRecursion(t *testing.T) {
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
	for i := 1; i <= 5; i++ {
		sum = sum + i;
	}
	sum;
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 15)
}

func TestBuiltinsAndStrings(t *testing.T) {
	input := `
	s := "hello";
	u := to_upper(s);
	len(u);
	`

	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 5)
}
