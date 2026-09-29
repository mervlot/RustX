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
var name string = "RustX";
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
		{"name", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
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
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("fn.Name is not 'add'. got=%s", fn.Name.Value)
	}
	if len(fn.Parameters) != 2 {
		t.Errorf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}
	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i);
}

for i < 5 {
	i++;
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements count wrong. got=%d", len(program.Statements))
	}

	for1, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if for1.Init == nil || for1.Condition == nil || for1.Post == nil {
		t.Errorf("3-part for loop components missing. Init=%v, Cond=%v, Post=%v", for1.Init, for1.Condition, for1.Post)
	}

	for2, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.ForStatement. got=%T", program.Statements[1])
	}
	if for2.Init != nil || for2.Condition == nil || for2.Post != nil {
		t.Errorf("single condition for loop incorrect. Init=%v, Cond=%v, Post=%v", for2.Init, for2.Condition, for2.Post)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}
