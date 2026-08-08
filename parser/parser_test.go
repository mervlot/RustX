package parser

import (
	"strings"
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y = 10;
var foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral() not '%s'. got=%s", name, varStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestForLoopParsing(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not expression statement. got=%T", program.Statements[0])
	}

	forExpr, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expression is not ForExpression. got=%T", stmt.Expression)
	}

	if forExpr.Init == nil {
		t.Fatalf("forExpr.Init is nil")
	}

	if forExpr.Condition == nil {
		t.Fatalf("forExpr.Condition is nil")
	}

	if forExpr.Post == nil {
		t.Fatalf("forExpr.Post is nil")
	}

	if len(forExpr.Body.Statements) != 1 {
		t.Fatalf("forExpr.Body.Statements has wrong length. got=%d", len(forExpr.Body.Statements))
	}
}

func TestFunctionParsing(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not expression statement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expression is not FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Fatalf("expected name 'add', got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type.Literal != "int" {
		t.Fatalf("parameter 0 is wrong")
	}

	if fn.ReturnType.Literal != "int" {
		t.Fatalf("expected return type int, got=%s", fn.ReturnType.Literal)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a * b + c",
			"((a * b) + c)",
		},
		{
			"a + b && c || d",
			"(((a + b) && c) || d)",
		},
		{
			"a == b != c",
			"((a == b) != c)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestSyntaxErrorLineAndColumn(t *testing.T) {
	input := `func add(a, b) {
	return a +
}`
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Fatalf("expected parser errors, got none")
	}

	firstErr := p.Errors()[0]
	if !strings.Contains(firstErr, "SyntaxError:") {
		t.Errorf("expected SyntaxError format, got %q", firstErr)
	}
	if !strings.Contains(firstErr, "line") {
		t.Errorf("expected error to contain line info, got %q", firstErr)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
