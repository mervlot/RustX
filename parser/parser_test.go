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

func TestShorthandAndAssignExpression(t *testing.T) {
	input := `
name := "RustX";
name = "RustX Starter";
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	assign1, ok := stmt1.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("stmt1.Expression is not ast.AssignExpression. got=%T", stmt1.Expression)
	}

	if assign1.Name.Value != "name" || assign1.Operator != ":=" {
		t.Errorf("assign1 wrong. got name=%s, op=%s", assign1.Name.Value, assign1.Operator)
	}

	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not ast.ExpressionStatement. got=%T", program.Statements[1])
	}

	assign2, ok := stmt2.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("stmt2.Expression is not ast.AssignExpression. got=%T", stmt2.Expression)
	}

	if assign2.Name.Value != "name" || assign2.Operator != "=" {
		t.Errorf("assign2 wrong. got name=%s, op=%s", assign2.Name.Value, assign2.Operator)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}
for x < 5 {
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

	forStmt1, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ForStatement. got=%T", program.Statements[0])
	}

	if forStmt1.Init == nil || forStmt1.Condition == nil || forStmt1.Post == nil {
		t.Fatalf("forStmt1 missing init/condition/post")
	}

	forStmt2, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[1] not *ast.ForStatement. got=%T", program.Statements[1])
	}

	if forStmt2.Init != nil || forStmt2.Condition == nil || forStmt2.Post != nil {
		t.Fatalf("forStmt2 wrong structure")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name != "add" {
		t.Errorf("function name wrong. expected='add', got=%q", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. expected=2, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("parameter 0 wrong")
	}

	if fn.ReturnType != "int" {
		t.Errorf("function return type wrong. expected='int', got=%q", fn.ReturnType)
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
