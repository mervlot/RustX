package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5
var y = 10
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", "int"},
		{"y", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}

		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}

		if tt.expectedType != "" {
			if varStmt.Type == nil || varStmt.Type.Value != tt.expectedType {
				t.Errorf("varStmt.Type wrong. expected=%s, got=%v", tt.expectedType, varStmt.Type)
			}
		}
	}
}

func TestColonAssignExpression(t *testing.T) {
	input := `name := "RustX"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.ColonAssignExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.ColonAssignExpression. got=%T", stmt.Expression)
	}

	if exp.Name.Value != "name" {
		t.Errorf("exp.Name.Value not 'name'. got=%s", exp.Name.Value)
	}

	strLit, ok := exp.Value.(*ast.StringLiteral)
	if !ok || strLit.Value != "RustX" {
		t.Errorf("exp.Value not 'RustX'. got=%v", exp.Value)
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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("fn.Name wrong. expected 'add', got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length wrong. expected 2, got=%d", len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Value != "int" {
		t.Errorf("fn.ReturnType wrong. expected 'int', got=%v", fn.ReturnType)
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

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForStatement. got=%T", program.Statements[0])
	}

	if forStmt.Init == nil {
		t.Errorf("forStmt.Init is nil")
	}
	if forStmt.Condition == nil {
		t.Errorf("forStmt.Condition is nil")
	}
	if forStmt.Post == nil {
		t.Errorf("forStmt.Post is nil")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %s", msg.Error())
	}
	t.FailNow()
}
