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

	return true
}

func TestAssignAndShorthand(t *testing.T) {
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

	assign1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok || !assign1.IsShorthand || assign1.Name.Value != "x" {
		t.Fatalf("statement[0] wrong format")
	}

	assign2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok || assign2.IsShorthand || assign2.Name.Value != "y" {
		t.Fatalf("statement[1] wrong format")
	}
}

func TestFunctionDeclarationAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
add(2, 3)
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok || fn.Name != "add" || len(fn.Parameters) != 2 || fn.ReturnType != "int" {
		t.Fatalf("function statement wrong format")
	}

	exprStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement[1] not ExpressionStatement")
	}

	call, ok := exprStmt.Expression.(*ast.CallExpression)
	if !ok || len(call.Arguments) != 2 {
		t.Fatalf("call expression wrong format")
	}
}

func TestForLoopAndWhile(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}
while x < 5 {
	x++
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
	if !ok || forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("for statement wrong format")
	}

	whileStmt, ok := program.Statements[1].(*ast.WhileStatement)
	if !ok || whileStmt.Condition == nil {
		t.Fatalf("while statement wrong format")
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
