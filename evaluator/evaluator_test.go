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
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.5", 5.5},
		{"10.0", 10.0},
		{"-5.5", -5.5},
		{"1.1 + 2.2", 3.3},
		{"5.0 * 2", 10.0},
		{"5 * 2.5", 12.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	// allow a small margin for float inaccuracy
	if result.Value < expected-0.0001 || result.Value > expected+0.0001 {
		t.Errorf("object has wrong value. got=%g, want=%g", result.Value, expected)
		return false
	}
	return true
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
		{"5.5 > 5.0", true},
		{"5.5 == 5.5", true},
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if 10 > 1 {
  if 10 > 1 {
    return 10;
  }
  return 1;
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

func TestVarAndShorthandAssign(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"a := 10; a;", 10},
		{"a := 10; a = 20; a;", 20},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{ (x + 2) }"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"func identity(x) { return x; }; identity(5);", 5},
		{"func identity(x) { x; }; identity(5);", 5},
		{"func double(x) { return x * 2; }; double(5);", 10},
		{"func add(x, y) { return x + y; }; add(5, 5);", 10},
		{"func add(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{`
func factorial(n) {
	if n <= 1 {
		return 1
	}
	return n * factorial(n - 1)
}
factorial(5)
`, 120},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"hello" == "hello"`, true},
		{`"hello" != "world"`, true},
		{`"hello" == "world"`, false},
	}

	for _, tt := range tests {
		testBooleanObject(t, testEval(tt.input), tt.expected)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"var i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"var myArray = [1, 2, 3]; myArray[2];", 3},
		{"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i];", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
		{`"hello"[0]`, "h"},
		{`"hello"[4]`, "o"},
		{`"hello"[5]`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("expected string, got=%T", evaluated)
			}
			if str.Value != expected {
				t.Errorf("expected %q, got=%q", expected, str.Value)
			}
		default:
			if evaluated.Type() != object.ERROR_OBJ {
				t.Errorf("expected error, got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}

func TestArrayIndexAssignment(t *testing.T) {
	input := `
var arr = [1, 2, 3]
arr[1] = 42
arr[1]
`
	testIntegerObject(t, testEval(input), 42)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`contains("hello", "ell")`, true},
		{`contains("hello", "world")`, false},
		{`to_upper("hello")`, "HELLO"},
		{`to_lower("WORLD")`, "world"},
		{`abs(-5)`, 5},
		{`abs(-5.5)`, 5.5},
		{`pow(2, 3)`, 8.0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case float64:
			testFloatObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if str.Value != expected {
				t.Errorf("string has wrong value. got=%q, want=%q", str.Value, expected)
			}
		}
	}
}

func TestDivisionByZero(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"5 / 0"},
		{"5.5 / 0.0"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated.Type() != object.ERROR_OBJ {
			t.Errorf("expected division-by-zero error, got=%T (%+v)", evaluated, evaluated)
		}
	}
}

func TestLogicalAndOrShortCircuit(t *testing.T) {
	// If left of && is false, right should not be evaluated (or cause error)
	inputAnd := `false && (5 / 0)`
	evaluatedAnd := testEval(inputAnd)
	testBooleanObject(t, evaluatedAnd, false)

	// If left of || is true, right should not be evaluated
	inputOr := `true || (5 / 0)`
	evaluatedOr := testEval(inputOr)
	testBooleanObject(t, evaluatedOr, true)
}

func TestPostfixIncrementDecrement(t *testing.T) {
	inputInc := `
x := 5
var original = x++
original
`
	testIntegerObject(t, testEval(inputInc), 5)

	inputIncValue := `
x := 5
x++
x
`
	testIntegerObject(t, testEval(inputIncValue), 6)
}

func TestWhileLoop(t *testing.T) {
	input := `
i := 0
while i < 5 {
	i++
}
i
`
	testIntegerObject(t, testEval(input), 5)
}

func TestForLoops(t *testing.T) {
	input := `
sum := 0
for i := 0; i < 5; i++ {
	sum = sum + i
}
sum
`
	testIntegerObject(t, testEval(input), 10)
}
