package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y int = 10;
var z string = "hello";
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors: %v", len(p.Errors()), p.Errors())
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", ""},
		{"y", "int"},
		{"z", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("test[%d] - statement not *ast.VarStatement. got=%T", i, stmt)
		}

		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("test[%d] - name value wrong. expected=%q, got=%q", i, tt.expectedIdentifier, varStmt.Name.Value)
		}

		if varStmt.Type != tt.expectedType {
			t.Fatalf("test[%d] - type wrong. expected=%q, got=%q", i, tt.expectedType, varStmt.Type)
		}
	}
}

func TestFunctionLiterals(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}
func() {
	return;
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors: %v", len(p.Errors()), p.Errors())
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	// First function: named add
	stmt0 := program.Statements[0].(*ast.ExpressionStatement)
	fn0, ok := stmt0.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("first statement is not *ast.FunctionLiteral. got=%T", stmt0.Expression)
	}

	if fn0.Name != "add" {
		t.Fatalf("fn0 name wrong. expected='add', got=%q", fn0.Name)
	}

	if len(fn0.Parameters) != 2 {
		t.Fatalf("fn0 should have 2 parameters. got=%d", len(fn0.Parameters))
	}

	if fn0.Parameters[0].Name.Value != "a" || fn0.Parameters[0].Type != "int" {
		t.Fatalf("fn0 parameter 0 wrong. got Name=%q, Type=%q", fn0.Parameters[0].Name.Value, fn0.Parameters[0].Type)
	}

	if fn0.ReturnType != "int" {
		t.Fatalf("fn0 ReturnType wrong. expected='int', got=%q", fn0.ReturnType)
	}

	// Second function: anonymous
	stmt1 := program.Statements[1].(*ast.ExpressionStatement)
	fn1, ok := stmt1.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("second statement is not *ast.FunctionLiteral. got=%T", stmt1.Expression)
	}

	if fn1.Name != "" {
		t.Fatalf("fn1 should be anonymous. got=%q", fn1.Name)
	}
}

func TestControlFlow(t *testing.T) {
	input := `
if x < y {
	println("yes");
} else if x == y {
	println("equal");
} else {
	println("no");
}

while flag {
	x++;
}

for i := 0; i < 10; i++ {
	println(i);
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		t.Fatalf("parser has %d errors: %v", len(p.Errors()), p.Errors())
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// Test If
	stmt0, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt0 not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	_, ok = stmt0.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr0 not *ast.IfExpression. got=%T", stmt0.Expression)
	}

	// Test While
	stmt1, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt1 not *ast.ExpressionStatement. got=%T", program.Statements[1])
	}
	_, ok = stmt1.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("expr1 not *ast.WhileExpression. got=%T", stmt1.Expression)
	}

	// Test For
	stmt2, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt2 not *ast.ExpressionStatement. got=%T", program.Statements[2])
	}
	forExpr, ok := stmt2.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expr2 not *ast.ForExpression. got=%T", stmt2.Expression)
	}

	if forExpr.Init == nil {
		t.Fatalf("forExpr.Init is nil")
	}
	if forExpr.Condition == nil {
		t.Fatalf("forExpr.Condition is nil")
	}
	if forExpr.Post == nil {
		t.Fatalf("forExpr.Post is nil")
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a || b && c",
			"(a || (b && c))",
		},
		{
			"a && b || c && d",
			"((a && b) || (c && d))",
		},
		{
			"x++",
			"(x++)",
		},
		{
			"x-- + y",
			"((x--) + y)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			t.Fatalf("parser error for input %q: %v", tt.input, p.Errors())
		}

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q for input %q", tt.expected, actual, tt.input)
		}
	}
}
