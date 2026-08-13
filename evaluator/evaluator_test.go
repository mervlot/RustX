package evaluator

import (
	"strings"
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
		{"12.34", 12.34},
		{"-1.5", -1.5},
		{"2.5 * 2.0", 5.0},
		{"5.0 / 2.0", 2.5},
		{"2.0 + 3", 5.0}, // coercion
		{"3 + 2.0", 5.0}, // coercion
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
		{"\"abc\" == \"abc\"", true},
		{"\"abc\" != \"def\"", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestShortCircuiting(t *testing.T) {
	// If short circuiting works, the right hand side with error or undefined identifier will not be evaluated!
	tests := []struct {
		input    string
		expected bool
	}{
		{"false && undefined_ident", false},
		{"true || undefined_ident", true},
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
		{"5.0 / 0.0", "division by zero"},
		{"5 / 0.0", "division by zero"},
		{"5.0 / 0", "division by zero"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Fatalf("no error returned for %s. got=%T (%+v)", tt.input, evaluated, evaluated)
		}
		if !strings.HasPrefix(errObj.Message, tt.expectedErr) {
			t.Errorf("wrong error message. expected prefix=%q, got=%q", tt.expectedErr, errObj.Message)
		}
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
		{"if 1 < 2 { 10 } else { 20 }", 10},
		{"if 1 > 2 { 10 } else { 20 }", 20},
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

func TestVarAndAssignments(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 5; a;", 25},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{"a := 10; a = 20; a;", 20},
		{"var x = 1; { x = 2; }; x;", 2}, // scoping update
		{"var x = 1; { x := 2; }; x;", 1}, // shadowing
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestWhileLoop(t *testing.T) {
	input := `
	var i = 0;
	while i < 5 {
		i = i + 1;
	}
	i;
	`
	testIntegerObject(t, testEval(input), 5)
}

func TestForLoop(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`
			var sum = 0;
			for i := 0; i < 5; i++ {
				sum = sum + i;
			}
			sum;
			`,
			10,
		},
		{
			`
			var sum = 0;
			var i = 0;
			for i < 3 {
				sum = sum + i;
				i = i + 1;
			}
			sum;
			`,
			3,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionLiteral(t *testing.T) {
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
		{"var identity = func(x) { x; }; identity(5);", 5},
		{"var identity = func(x) { return x; }; identity(5);", 5},
		{"var double = func(x) { x * 2; }; double(5);", 10},
		{"var add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
		{
			`
			func fib(n) {
				if n <= 1 { return n; }
				return fib(n - 1) + fib(n - 2);
			}
			fib(6);
			`,
			8,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestArrayLiteralsAndIndexing(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][2]", 3},
		{"var a = [1, 2, 3]; a[1] = 99; a[1]", 99},
		{"len([1, 2, 3])", 3},
		{"\"hello\"[1]", "e"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if integer, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestStringBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"contains(\"hello world\", \"world\")", true},
		{"to_upper(\"abc\")", "ABC"},
		{"to_lower(\"XYZ\")", "xyz"},
		{"join([\"a\", \"b\", \"c\"], \"-\")", "a-b-c"},
		{"split(\"a-b-c\", \"-\")[1]", "b"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if b, ok := tt.expected.(bool); ok {
			testBooleanObject(t, evaluated, b)
		} else {
			testStringObject(t, evaluated, tt.expected.(string))
		}
	}
}

func TestMathBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"abs(-12.5)", 12.5},
		{"pow(2, 3)", 8.0},
	}

	for _, tt := range tests {
		testFloatObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFilesystemBuiltins(t *testing.T) {
	input := `
	write_file("test_temp.txt", "rustx file test content");
	var s = read_file("test_temp.txt");
	s;
	`
	testStringObject(t, testEval(input), "rustx file test content")
}

// Helpers

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

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q", result.Value, expected)
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
