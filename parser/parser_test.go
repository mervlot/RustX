package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5
var y int = 10
var name = "rustx"
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
		{"name", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not *ast.VarStatement. got=%T", i, stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not %q. got=%q", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Fatalf("varStmt.Type not %q. got=%q", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShortVarStatement(t *testing.T) {
	input := `age := 25`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}
	if stmt.Name.Value != "age" {
		t.Fatalf("stmt.Name.Value not 'age'. got=%q", stmt.Name.Value)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}
	if fn.Name != "add" {
		t.Fatalf("fn.Name not 'add'. got=%q", fn.Name)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn parameters count wrong. got=%d", len(fn.Parameters))
	}
	if fn.ReturnType != "int" {
		t.Fatalf("fn return type wrong. got=%q", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i) }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExpr, ok := exprStmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expr is not *ast.ForExpression. got=%T", exprStmt.Expression)
	}
	if forExpr.Init == nil || forExpr.Condition == nil || forExpr.Post == nil {
		t.Fatalf("forExpr components nil: init=%v, cond=%v, post=%v", forExpr.Init, forExpr.Condition, forExpr.Post)
	}
}

func TestIfElseParsing(t *testing.T) {
	input := `if x > 10 { y = 1 } else { y = 2 }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr is not *ast.IfExpression. got=%T", exprStmt.Expression)
	}
	if ifExpr.Alternative == nil {
		t.Fatalf("ifExpr.Alternative is nil")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
