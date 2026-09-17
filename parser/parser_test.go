package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarAndAssignStatements(t *testing.T) {
	input := `
var x int = 5;
var y = 10;
name := "RustX";
x = 15;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	// 1. var x int = 5
	varStmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt[0] is not *ast.VarStatement. got=%T", program.Statements[0])
	}
	if varStmt.Name.Value != "x" || varStmt.Type != "int" {
		t.Errorf("varStmt.Name=%s, Type=%s", varStmt.Name.Value, varStmt.Type)
	}

	// 2. var y = 10
	varStmt2, ok := program.Statements[1].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt[1] is not *ast.VarStatement. got=%T", program.Statements[1])
	}
	if varStmt2.Name.Value != "y" || varStmt2.Type != "" {
		t.Errorf("varStmt2.Name=%s, Type=%s", varStmt2.Name.Value, varStmt2.Type)
	}

	// 3. name := "RustX"
	assignStmt, ok := program.Statements[2].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt[2] is not *ast.AssignStatement. got=%T", program.Statements[2])
	}
	if assignStmt.Name.Value != "name" || assignStmt.Token.Literal != ":=" {
		t.Errorf("assignStmt.Name=%s, Token=%s", assignStmt.Name.Value, assignStmt.Token.Literal)
	}

	// 4. x = 15
	assignStmt2, ok := program.Statements[3].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt[3] is not *ast.AssignStatement. got=%T", program.Statements[3])
	}
	if assignStmt2.Name.Value != "x" || assignStmt2.Token.Literal != "=" {
		t.Errorf("assignStmt2.Name=%s, Token=%s", assignStmt2.Name.Value, assignStmt2.Token.Literal)
	}
}

func TestFunctionDeclaration(t *testing.T) {
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
		t.Fatalf("stmt[0] is not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fn.Name != "add" {
		t.Errorf("fn.Name is not 'add'. got=%q", fn.Name)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}
	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("param 0 wrong")
	}
	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType wrong. got=%q", fn.ReturnType)
	}
}

func TestForLoopVariants(t *testing.T) {
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
		t.Fatalf("program.Statements count wrong. got=%d", len(program.Statements))
	}

	// 1. 3-part for loop
	forStmt1, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[0] is not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if forStmt1.Init == nil || forStmt1.Condition == nil || forStmt1.Post == nil {
		t.Fatalf("forStmt1 missing parts")
	}

	// 2. Single condition for loop
	forStmt2, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[1] is not *ast.ForStatement. got=%T", program.Statements[1])
	}
	if forStmt2.Init != nil || forStmt2.Condition == nil {
		t.Fatalf("forStmt2 parts wrong")
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
