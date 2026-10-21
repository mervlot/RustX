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
name := "RustX"
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
		{"name"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		switch s := stmt.(type) {
		case *ast.VarStatement:
			if s.Name.Value != tt.expectedIdentifier {
				t.Errorf("s.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, s.Name.Value)
			}
		case *ast.ShortVarStatement:
			if s.Name.Value != tt.expectedIdentifier {
				t.Errorf("s.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, s.Name.Value)
			}
		default:
			t.Fatalf("stmt not VarStatement or ShortVarStatement. got=%T", stmt)
		}
	}
}

func TestFunctionDeclaration(t *testing.T) {
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

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr is not FunctionLiteral. got=%T", exprStmt.Expression)
	}

	if fn.Name == nil || fn.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%v", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. expected 2, got=%d", len(fn.Parameters))
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
		t.Fatalf("stmt is not ForStatement. got=%T", program.Statements[0])
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

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %s", msg.Error())
	}
	t.FailNow()
}
