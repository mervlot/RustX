package parser

import (
	"testing"

	"rustx/ast"
	"rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5;
var y = 10;
name := "Muhammad";
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if err := CheckParserErrors(p); err != nil {
		t.Fatalf("parser errors: %v", err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// 1. var x int = 5
	varStmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("s is not *ast.VarStatement. got=%T", program.Statements[0])
	}
	if varStmt.Name.Value != "x" {
		t.Fatalf("varStmt.Name.Value not 'x'. got=%s", varStmt.Name.Value)
	}
	if varStmt.Type != "int" {
		t.Fatalf("varStmt.Type not 'int'. got=%s", varStmt.Type)
	}

	// 2. var y = 10
	varStmt2, ok := program.Statements[1].(*ast.VarStatement)
	if !ok {
		t.Fatalf("s is not *ast.VarStatement. got=%T", program.Statements[1])
	}
	if varStmt2.Name.Value != "y" {
		t.Fatalf("varStmt2.Name.Value not 'y'. got=%s", varStmt2.Name.Value)
	}

	// 3. name := "Muhammad"
	assignStmt, ok := program.Statements[2].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("s is not *ast.AssignStatement. got=%T", program.Statements[2])
	}
	if assignStmt.Name.Value != "name" {
		t.Fatalf("assignStmt.Name.Value not 'name'. got=%s", assignStmt.Name.Value)
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

	if err := CheckParserErrors(p); err != nil {
		t.Fatalf("parser errors: %v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fnStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 0 not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := fnStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr not *ast.FunctionLiteral. got=%T", fnStmt.Expression)
	}

	if fn.Name != "add" {
		t.Fatalf("fn.Name not 'add'. got=%s", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType != "int" {
		t.Fatalf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i);
}

for i < 5 {
	i = i + 1;
}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if err := CheckParserErrors(p); err != nil {
		t.Fatalf("parser errors: %v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements count wrong. got=%d", len(program.Statements))
	}

	forStmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 0 not ExpressionStatement. got=%T", program.Statements[0])
	}

	forExpr1, ok := forStmt1.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expr 0 not ForExpression. got=%T", forStmt1.Expression)
	}

	if forExpr1.Init == nil || forExpr1.Condition == nil || forExpr1.Post == nil {
		t.Fatalf("3-part for loop component missing: init=%v, cond=%v, post=%v",
			forExpr1.Init, forExpr1.Condition, forExpr1.Post)
	}

	forStmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 1 not ExpressionStatement. got=%T", program.Statements[1])
	}

	forExpr2, ok := forStmt2.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expr 1 not ForExpression. got=%T", forStmt2.Expression)
	}

	if forExpr2.Condition == nil {
		t.Fatalf("single condition for loop missing condition")
	}
}
