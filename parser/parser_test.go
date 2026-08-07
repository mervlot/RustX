package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y = 10.5;
var foobar = "hello";
z := true;
`
	l := lexer.New(input)
	p := New(l, "test.rx")
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
		{"z"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got='%s'", name, varStmt.Name.Value)
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
	p := New(l, "test.rx")
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be ExpressionStatement, got=%T", program.Statements[0])
	}

	forExpr, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("Expected Expression to be *ast.ForExpression, got=%T", stmt.Expression)
	}

	if forExpr.Init == nil {
		t.Errorf("Expected Init not nil")
	}
	if forExpr.Condition == nil {
		t.Errorf("Expected Condition not nil")
	}
	if forExpr.Post == nil {
		t.Errorf("Expected Post not nil")
	}
}

func TestFunctionLiteralWithTypes(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
`
	l := lexer.New(input)
	p := New(l, "test.rx")
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected ExpressionStatement, got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expected FunctionLiteral, got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("Expected fn.Name.Value 'add', got=%q", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("Expected 2 params, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Type != "int" {
		t.Errorf("Expected first param type 'int', got=%q", fn.Parameters[0].Type)
	}

	if fn.ReturnType != "int" {
		t.Errorf("Expected return type 'int', got=%q", fn.ReturnType)
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
