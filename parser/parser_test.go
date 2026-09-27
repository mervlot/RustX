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
		{"x", "int"},
		{"y", ""},
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

func TestShorthandAndAssignmentStatements(t *testing.T) {
	input := `
x := 5;
y = 10;
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
		t.Fatalf("stmt[0] not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}
	if shortStmt.Name.Value != "x" {
		t.Errorf("shortStmt.Name.Value not 'x'. got=%s", shortStmt.Name.Value)
	}

	assignStmt, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.AssignmentStatement. got=%T", program.Statements[1])
	}
	if assignStmt.Name.Value != "y" {
		t.Errorf("assignStmt.Name.Value not 'y'. got=%s", assignStmt.Name.Value)
	}
}

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}

add(5, 10);
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}

	if fn.Name != "add" {
		t.Errorf("fn.Name not 'add'. got=%s", fn.Name)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length not 2. got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("param 0 mismatch. got=%v %v", fn.Parameters[0].Name.Value, fn.Parameters[0].Type)
	}
	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForAndWhileLoopParsing(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i);
}

while x < 5 {
	x++;
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements length not 2. got=%d", len(program.Statements))
	}

	_, okFor := program.Statements[0].(*ast.ForStatement)
	if !okFor {
		t.Fatalf("stmt[0] not *ast.ForStatement. got=%T", program.Statements[0])
	}

	_, okWhile := program.Statements[1].(*ast.WhileStatement)
	if !okWhile {
		t.Fatalf("stmt[1] not *ast.WhileStatement. got=%T", program.Statements[1])
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
