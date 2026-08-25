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
var foobar = 838383;
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

func TestAssignmentStatements(t *testing.T) {
	input := `
x := 5
y = 10
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	stmt0, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt0 not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}
	if stmt0.Name.Value != "x" || stmt0.Token.Literal != ":=" {
		t.Errorf("stmt0 mismatch. got name=%s token=%s", stmt0.Name.Value, stmt0.Token.Literal)
	}

	stmt1, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt1 not *ast.AssignmentStatement. got=%T", program.Statements[1])
	}
	if stmt1.Name.Value != "y" || stmt1.Token.Literal != "=" {
		t.Errorf("stmt1 mismatch. got name=%s token=%s", stmt1.Name.Value, stmt1.Token.Literal)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(x int, y int) int { return x + y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if function.Name == nil || function.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%v", function.Name)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	if function.ReturnType == nil || function.ReturnType.Value != "int" {
		t.Fatalf("function return type is not 'int'. got=%v", function.ReturnType)
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
		t.Fatalf("stmt is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("forStmt init/condition/post missing")
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
