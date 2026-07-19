package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y int = 10;
var foobar float = 12.34;
`
	l := lexer.New(input, "test.rx")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", ""},
		{"y", "int"},
		{"foobar", "float"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, typeAnn string) bool {
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got='%s'", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Type != typeAnn {
		t.Errorf("varStmt.Type not '%s'. got='%s'", typeAnn, varStmt.Type)
		return false
	}

	return true
}

func TestShorthandAndAssign(t *testing.T) {
	input := `
x := 5
x = 10
`
	l := lexer.New(input, "test.rx")
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	// 1. Shorthand
	s1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 0 is not expression statement. got=%T", program.Statements[0])
	}
	assign1, ok := s1.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("exp 0 is not AssignExpression. got=%T", s1.Expression)
	}
	if assign1.Name.Value != "x" || assign1.Operator != ":=" {
		t.Errorf("expected shorthand x := 5, got %s %s", assign1.Name.Value, assign1.Operator)
	}

	// 2. Assignment
	s2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 1 is not expression statement. got=%T", program.Statements[1])
	}
	assign2, ok := s2.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("exp 1 is not AssignExpression. got=%T", s2.Expression)
	}
	if assign2.Name.Value != "x" || assign2.Operator != "=" {
		t.Errorf("expected assignment x = 10, got %s %s", assign2.Name.Value, assign2.Operator)
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b * c", "(a + (b * c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; 5 * 5", "(3 + 4)(5 * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 || 2 < 4 && 5 == 5", "((3 > 5) || ((2 < 4) && (5 == 5)))"},
		{"a * [1, 2, 3][b]", "(a * ([1, 2, 3][b]))"},
		{"add(a, b) + c", "(add(a, b) + c)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input, "test.rx")
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if x < y { x } else { y }`

	l := lexer.New(input, "test.rx")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if exp.Condition.String() != "(x < y)" {
		t.Errorf("condition was not (x < y). got=%s", exp.Condition.String())
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("consequence statement is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if consequence.Expression.String() != "x" {
		t.Errorf("consequence statement expression was not x. got=%s", consequence.Expression.String())
	}

	if exp.Alternative == nil {
		t.Fatalf("exp.Alternative was nil")
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got=%d", len(exp.Alternative.Statements))
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("alternative statement is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if alt.Expression.String() != "y" {
		t.Errorf("alternative statement expression was not y. got=%s", alt.Expression.String())
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(x int, y int) int { return x + y; }`

	l := lexer.New(input, "test.rx")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if function.Name != "add" {
		t.Errorf("function name is not 'add'. got='%s'", function.Name)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	if function.Parameters[0].Name.Value != "x" || function.Parameters[0].Type != "int" {
		t.Errorf("parameter 0 is wrong. got=%+v", function.Parameters[0])
	}

	if function.ReturnType != "int" {
		t.Errorf("return type is wrong. expected 'int', got=%q", function.ReturnType)
	}
}

func TestLoops(t *testing.T) {
	input := `
while x < 10 {
	x = x + 1
}

for i := 0; i < 10; i++ {
	println(i)
}
`
	l := lexer.New(input, "test.rx")
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	// 1. While
	s1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 0 is not expression statement. got=%T", program.Statements[0])
	}
	whileExp, ok := s1.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("exp 0 is not WhileExpression. got=%T", s1.Expression)
	}
	if whileExp.Condition.String() != "(x < 10)" {
		t.Errorf("expected while condition to be (x < 10), got=%s", whileExp.Condition.String())
	}

	// 2. For
	s2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 1 is not expression statement. got=%T", program.Statements[1])
	}
	forExp, ok := s2.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("exp 1 is not ForExpression. got=%T", s2.Expression)
	}
	if forExp.Init.String() != "(i := 0)" {
		t.Errorf("expected init to be (i := 0), got=%s", forExp.Init.String())
	}
	if forExp.Condition.String() != "(i < 10)" {
		t.Errorf("expected condition to be (i < 10), got=%s", forExp.Condition.String())
	}
	if forExp.Post.String() != "(i++)" {
		t.Errorf("expected post to be (i++), got=%s", forExp.Post.String())
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, err := range errs {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
