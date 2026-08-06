package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5;
var y = 10;
var z string;
`
	l := lexer.New(input)
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
		{"x", "int"},
		{"y", ""},
		{"z", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, expectedType string) bool {
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got='%s'", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Type != expectedType {
		t.Errorf("varStmt.Type not '%s'. got='%s'", expectedType, varStmt.Type)
		return false
	}

	return true
}

func TestShorthandStatements(t *testing.T) {
	input := `x := 5;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ShorthandStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ShorthandStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%q", stmt.Name.Value)
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a || b && c",
			"(a || (b && c))",
		},
		{
			"a && b || c",
			"((a && b) || c)",
		},
		{
			"a + b * c",
			"(a + (b * c))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected+";" {
			t.Errorf("expected=%q, got=%q", tt.expected+";", actual)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(x int, y int) int { return x + y; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name != "add" {
		t.Errorf("expected name 'add', got %q", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "x" || fn.Parameters[0].Type != "int" {
		t.Errorf("param 0 wrong. got=%q type=%q", fn.Parameters[0].Name.Value, fn.Parameters[0].Type)
	}

	if fn.ReturnType != "int" {
		t.Errorf("expected return type 'int', got %q", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i); }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.ForExpression. got=%T", stmt.Expression)
	}

	if forExp.Init == nil {
		t.Fatalf("forExp.Init is nil")
	}

	shorthand, ok := forExp.Init.(*ast.ShorthandStatement)
	if !ok {
		t.Fatalf("forExp.Init not *ast.ShorthandStatement. got=%T", forExp.Init)
	}
	if shorthand.Name.Value != "i" {
		t.Errorf("expected init identifier 'i', got %q", shorthand.Name.Value)
	}

	if forExp.Condition == nil {
		t.Fatalf("forExp.Condition is nil")
	}

	if forExp.Post == nil {
		t.Fatalf("forExp.Post is nil")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
