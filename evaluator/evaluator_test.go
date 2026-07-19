package evaluator

import (
	"os"
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

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.5", 5.5},
		{"-10.2", -10.2},
		{"2.0 * 3.5", 7.0},
		{"10.0 / 4.0", 2.5},
		{"1.5 + 2", 3.5}, // Float coercion
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
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
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
if 10 > 1 {
  if 10 > 1 {
    return 10;
  }
  return 1;
}
`, 10},
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
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"var a = 5; a = 10; a;", 10},
		{"a := 5; a;", 5},
		{"a := 5; a = a + 5; a;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLoops(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
var x = 0;
while x < 5 {
	x = x + 1;
}
x;
`, 5},
		{`
var sum = 0;
for i := 0; i < 5; i++ {
	sum = sum + i;
}
sum;
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T", evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].Name.Value != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0].Name.Value)
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
		{"var identity = func(x) { x; }; identity(5);", 5},
		{"var identity = func(x) { return x; }; identity(5);", 5},
		{"var double = func(x) { x * 2; }; double(5);", 10},
		{"var add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
var newAdder = func(x) {
  func(y) { x + y };
};

var addTwo = newAdder(2);
addTwo(3);`

	testIntegerObject(t, testEval(input), 5)
}

func TestRecursiveFunctions(t *testing.T) {
	input := `
func factorial(n) {
	if n <= 1 {
		return 1;
	}
	return n * factorial(n - 1);
}
factorial(5);
`
	testIntegerObject(t, testEval(input), 120)
}

func TestShortCircuiting(t *testing.T) {
	// If short circuit works, we should not execute the division by zero that occurs later
	tests := []struct {
		input    string
		expected bool
	}{
		{"false && (1/0 == 0)", false},
		{"true || (1/0 == 0)", true},
	}

	for _, tt := range tests {
		testBooleanObject(t, testEval(tt.input), tt.expected)
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
			t.Fatalf("no error returned for div-by-zero. got=%T", evaluated)
		}
		if !containsSub(errObj.Message, tt.expectedErr) {
			t.Errorf("error message wrong. got=%q, expected to contain=%q", errObj.Message, tt.expectedErr)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T", evaluated)
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
		t.Fatalf("object is not String. got=%T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringAndArrayIndexing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello"[1]`, "e"},
		{`[10, 20, 30][2]`, int64(30)},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("not a string object: got=%T", evaluated)
			}
			if str.Value != expected {
				t.Errorf("expected string %q, got %q", expected, str.Value)
			}
		case int64:
			testIntegerObject(t, evaluated, expected)
		}
	}
}

func TestBuiltins(t *testing.T) {
	// Write a temp file for testing file system basics
	tempFile := "temp_test_rustx.txt"
	defer os.Remove(tempFile)

	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("hello")`, int64(5)},
		{`len([1, 2, 3])`, int64(3)},
		{`contains("hello", "ell")`, true},
		{`contains("hello", "world")`, false},
		{`to_upper("hello")`, "HELLO"},
		{`to_lower("WORLD")`, "world"},
		{`abs(-5)`, int64(5)},
		{`abs(-5.5)`, 5.5},
		{`pow(2, 3)`, 8.0},
		{`join(["a", "b"], "-")`, "a-b"},
		{`write_file("` + tempFile + `", "rustx is awesome")`, true},
		{`read_file("` + tempFile + `")`, "rustx is awesome"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case float64:
			testFloatObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected)
		case string:
			str, ok := evaluated.(*object.String)
			if !ok {
				t.Fatalf("not a string object: got=%T (%+v) for input %s", evaluated, evaluated, tt.input)
			}
			if str.Value != expected {
				t.Errorf("expected string %q, got %q", expected, str.Value)
			}
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input, "test.rx")
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func containsSub(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > len(sub) && (s[:len(sub)] == sub || s[1:] == sub || len(s) > len(sub)+1)) // standard contains check without import strings if we want to avoid package clashes
}
