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

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"1.5 + 2.5", 4.0},
		{"5.0 - 2.5", 2.5},
		{"2.0 * 3.5", 7.0},
		{"7.0 / 2.0", 3.5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		flt, ok := evaluated.(*object.Float)
		if !ok {
			t.Fatalf("object is not Float. got=%T (%+v)", evaluated, evaluated)
		}
		if flt.Value != tt.expected {
			t.Errorf("object has wrong value. got=%f, want=%f", flt.Value, tt.expected)
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

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].Name.String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
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
			func factorial(n int) int {
				if n <= 1 {
					return 1
				}
				return n * factorial(n - 1)
			}
			factorial(5)
			`,
			120,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestWhileAndForLoops(t *testing.T) {
	inputWhile := `
	sum := 0
	i := 1
	while i <= 5 {
		sum = sum + i
		i++
	}
	sum
	`
	testIntegerObject(t, testEval(inputWhile), 15)

	inputFor := `
	sum := 0
	for i := 1; i <= 5; i++ {
		sum = sum + i
	}
	sum
	`
	testIntegerObject(t, testEval(inputFor), 15)
}

func TestBuiltins(t *testing.T) {
	inputStr := `
	s := "hello world"
	len(s)
	`
	testIntegerObject(t, testEval(inputStr), 11)

	inputUpper := `
	str_to_upper("rustx")
	`
	evaluated := testEval(inputUpper)
	str, ok := evaluated.(*object.String)
	if !ok || str.Value != "RUSTX" {
		t.Fatalf("expected RUSTX, got=%+v", evaluated)
	}

	inputMath := `
	math_sqrt(16)
	`
	evaluatedMath := testEval(inputMath)
	flt, ok := evaluatedMath.(*object.Float)
	if !ok || flt.Value != 4.0 {
		t.Fatalf("expected 4.0, got=%+v", evaluatedMath)
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
