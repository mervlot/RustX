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
z := 15
name := "RustX"
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		isShorthand        bool
	}{
		{"x", false},
		{"y", false},
		{"z", true},
		{"name", true},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("statements[%d] is not *ast.VarStatement. got=%T", i, stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.IsShorthand != tt.isShorthand {
			t.Errorf("varStmt.IsShorthand not %v. got=%v", tt.isShorthand, varStmt.IsShorthand)
		}
	}
}

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
result := add(5, 10)
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
		t.Fatalf("statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := fnStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("fnStmt.Expression is not *ast.FunctionLiteral. got=%T", fnStmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("function name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count not 2. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType != "int" {
		t.Errorf("function return type not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForAndWhileLoops(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}

while x > 0 {
	x = x - 1
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	_, ok = forStmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("forStmt.Expression is not *ast.ForExpression. got=%T", forStmt.Expression)
	}

	whileStmt, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statements[1] is not *ast.ExpressionStatement. got=%T", program.Statements[1])
	}

	_, ok = whileStmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("whileStmt.Expression is not *ast.WhileExpression. got=%T", whileStmt.Expression)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
