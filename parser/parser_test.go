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
		expectedType       string
	}{
		{"x", ""},
		{"y", "int"},
		{"foobar", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, typeAnn string) bool {
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

	if varStmt.Type != typeAnn {
		t.Errorf("varStmt.Type not '%s'. got=%s", typeAnn, varStmt.Type)
		return false
	}

	return true
}

func TestShortVarAndAssignmentStatements(t *testing.T) {
	input := `
x := 5;
x = 10;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	shortStmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}
	if shortStmt.Name.Value != "x" {
		t.Errorf("shortStmt.Name.Value not 'x'. got=%s", shortStmt.Name.Value)
	}

	assignStmt, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not *ast.AssignmentStatement. got=%T", program.Statements[1])
	}
	if assignStmt.Name.Value != "x" {
		t.Errorf("assignStmt.Name.Value not 'x'. got=%s", assignStmt.Name.Value)
	}
}

func TestControlFlowParsing(t *testing.T) {
	input := `
if x < y {
	return true
} else {
	return false
}

while x < 10 {
	x++
}

for i := 0; i < 10; i++ {
	println(i)
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
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
