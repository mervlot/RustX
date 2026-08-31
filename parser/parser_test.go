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
z := 15
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
		{"z"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
	}
}

func TestForLoop(t *testing.T) {
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
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ForStatement. got=%T", program.Statements[0])
	}

	if forStmt.Init == nil {
		t.Fatalf("forStmt.Init is nil")
	}
	if forStmt.Condition == nil {
		t.Fatalf("forStmt.Condition is nil")
	}
	if forStmt.Post == nil {
		t.Fatalf("forStmt.Post is nil")
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `
func add(a int, b int) int {
    return a + b
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
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
		t.Fatalf("fn.Parameters length not 2. got=%d", len(fn.Parameters))
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
