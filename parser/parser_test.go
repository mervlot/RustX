package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5
var y int = 10
var foobar = 838383
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
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
	}
}

func TestDefineAndAssignStatements(t *testing.T) {
	input := `
x := 10
x = 20
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	defStmt, ok := program.Statements[0].(*ast.DefineStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.DefineStatement. got=%T", program.Statements[0])
	}
	if defStmt.Name.Value != "x" {
		t.Errorf("defStmt.Name.Value not 'x'. got=%s", defStmt.Name.Value)
	}

	assignStmt, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.AssignmentStatement. got=%T", program.Statements[1])
	}
	if assignStmt.Name.Value != "x" {
		t.Errorf("assignStmt.Name.Value not 'x'. got=%s", assignStmt.Name.Value)
	}
}

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
add(5, 10)
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fnStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := fnStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("fnStmt.Expression not *ast.FunctionLiteral. got=%T", fnStmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("fn.Name.Value not 'add'. got=%s", fn.Name.Value)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length not 2. got=%d", len(fn.Parameters))
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
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ForStatement. got=%T", program.Statements[0])
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
