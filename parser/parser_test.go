package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y int = 10;
var foobar = true;
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
		{"x", ""},
		{"y", "int"},
		{"foobar", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] not *ast.VarStatement. got=%T", i, stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got='%s'", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShortVarStatements(t *testing.T) {
	input := `
name := "Muhammad";
age := 16;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got=%d", len(program.Statements))
	}

	shortStmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}
	if shortStmt.Name.Value != "name" {
		t.Errorf("shortStmt.Name not 'name'. got=%s", shortStmt.Name.Value)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
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
		t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name != "add" {
		t.Errorf("fn.Name not 'add'. got=%s", fn.Name)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("fn.Parameters[0] wrong: %v", fn.Parameters[0])
	}
	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
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
		t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.ForExpression. got=%T", stmt.Expression)
	}

	if forExp.Init == nil {
		t.Fatalf("forExp.Init is nil")
	}
	if forExp.Condition == nil {
		t.Fatalf("forExp.Condition is nil")
	}
	if forExp.Post == nil {
		t.Fatalf("forExp.Post is nil")
	}
}

func TestIfExpression(t *testing.T) {
	input := `
if x < y {
	x
} else {
	y
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
		t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression not *ast.IfExpression. got=%T", stmt.Expression)
	}

	if ifExp.Alternative == nil {
		t.Fatalf("ifExp.Alternative is nil")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %v", err)
	}
	t.FailNow()
}
