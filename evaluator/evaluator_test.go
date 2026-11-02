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
	prog := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(prog, env)
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
		{"true && false", false},
		{"true && true", true},
		{"false || true", true},
		{"false || false", false},
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
		{"if 1 > 2 { 10 }", nil},
		{"if 1 > 2 { 10 } else { 20 }", 20},
		{"if 1 < 2 { 10 } else { 20 }", 10},
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

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10", 10},
		{"return 10; 9", 10},
		{"return 2 * 5; 9", 10},
		{"9; return 2 * 5; 9", 10},
		{
			`
if 10 > 1 {
	if 10 > 1 {
		return 10
	}
	return 1
}
`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestVarAndAssign(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a", 5},
		{"var a = 5 * 5; a", 25},
		{"var a = 5; var b = a; b", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c", 15},
		{"a := 10; a", 10},
		{"a := 10; a = 20; a", 20},
		{"a := 10; a++; a", 11},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x int) { x + 2 }"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x int" {
		t.Fatalf("parameter is not 'x int'. got=%q", fn.Parameters[0])
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"func identity(x int) int { return x; }; identity(5)", 5},
		{"func identity(x int) int { x; }; identity(5)", 5},
		{"func double(x int) int { x * 2; }; double(5)", 10},
		{"func add(x int, y int) int { x + y; }; add(5, 5)", 10},
		{"func add(x int, y int) int { x + y; }; add(5 + 5, add(5, 5))", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
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

func TestWhileLoop(t *testing.T) {
	input := `
i := 0
sum := 0
while i < 5 {
	i++
	sum = sum + i
}
sum
`
	testIntegerObject(t, testEval(input), 15)
}

func TestForLoop(t *testing.T) {
	input := `
sum := 0
for i := 1; i <= 5; i++ {
	sum = sum + i
}
sum
`
	testIntegerObject(t, testEval(input), 15)
}

func TestBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, int64(5)},
		{`len([1, 2, 3])`, int64(3)},
		{`str_upper("hello")`, "HELLO"},
		{`str_lower("WORLD")`, "world"},
		{`str_contains("rustx starter", "start")`, true},
		{`math_sqrt(16)`, float64(4)},
		{`math_pow(2, 3)`, float64(8)},
		{`math_abs(-42)`, int64(42)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case float64:
			fl, ok := evaluated.(*object.Float)
			if !ok || fl.Value != expected {
				t.Fatalf("expected float %g, got %v", expected, evaluated)
			}
		case string:
			str, ok := evaluated.(*object.String)
			if !ok || str.Value != expected {
				t.Fatalf("expected string %s, got %v", expected, evaluated)
			}
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
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
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
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
