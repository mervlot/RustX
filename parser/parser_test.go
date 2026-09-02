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
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, expectedType string) bool {
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

	if varStmt.Type != expectedType {
		t.Errorf("varStmt.Type not '%s'. got=%s", expectedType, varStmt.Type)
		return false
	}

	return true
}

func TestShorthandAndAssignmentStatements(t *testing.T) {
	input := `
x := 10;
x = 20;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	decl, ok := program.Statements[0].(*ast.DeclAssignStatement)
	if !ok {
		t.Fatalf("stmt[0] not DeclAssignStatement. got=%T", program.Statements[0])
	}
	if decl.Name.Value != "x" {
		t.Errorf("decl.Name.Value not 'x'. got=%s", decl.Name.Value)
	}

	assign, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt[1] not AssignStatement. got=%T", program.Statements[1])
	}
	if assign.Name.Value != "x" {
		t.Errorf("assign.Name.Value not 'x'. got=%s", assign.Name.Value)
	}
}

func TestFunctionStatement(t *testing.T) {
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
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt not FunctionLiteral. got=%T", program.Statements[0])
	}

	if fn.Name != "add" {
		t.Errorf("fn.Name not 'add'. got=%s", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length wrong. got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("param 0 wrong. got=%v %v", fn.Parameters[0].Name.Value, fn.Parameters[0].Type)
	}

	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForAndWhileLoops(t *testing.T) {
	input := `
while x < 10 {
	x++;
}

for i := 0; i < 10; i++ {
	println(i);
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	_, okWhile := program.Statements[0].(*ast.WhileStatement)
	if !okWhile {
		t.Fatalf("stmt[0] not WhileStatement. got=%T", program.Statements[0])
	}

	forStmt, okFor := program.Statements[1].(*ast.ForStatement)
	if !okFor {
		t.Fatalf("stmt[1] not ForStatement. got=%T", program.Statements[1])
	}

	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Errorf("ForStatement components missing: init=%v, cond=%v, post=%v", forStmt.Init, forStmt.Condition, forStmt.Post)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, err := range errs {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
