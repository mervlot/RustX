package parser

import (
	"testing"

	"rustx/ast"
	"rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5
var y int = 10
var name string = "RustX"
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
		{"name", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Fatalf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
	}
}

func TestAssignStatements(t *testing.T) {
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

	s1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok || !s1.IsDefine || s1.Target.String() != "x" {
		t.Fatalf("statement 0 invalid assign statement")
	}

	s2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok || s2.IsDefine || s2.Target.String() != "y" {
		t.Fatalf("statement 1 invalid assign statement")
	}
}

func TestForAndWhileStatements(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
    x = x + i
}
while x > 0 {
    x--
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt 0 not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("forStmt header parts missing")
	}

	whileStmt, ok := program.Statements[1].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt 1 not *ast.WhileStatement. got=%T", program.Statements[1])
	}
	if whileStmt.Condition == nil {
		t.Fatalf("whileStmt condition missing")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
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
		t.Fatalf("expr not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name != "add" {
		t.Fatalf("fn.Name not 'add'. got=%s", fn.Name)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length not 2. got=%d", len(fn.Parameters))
	}
	if fn.ReturnType != "int" {
		t.Fatalf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
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
