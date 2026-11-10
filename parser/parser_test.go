package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5
var y = 10.5
var foobar = "hello"
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
		hasType            bool
		expectedType       string
	}{
		{"x", true, "int"},
		{"y", false, ""},
		{"foobar", false, ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if tt.hasType {
			if varStmt.Type == nil || varStmt.Type.Name != tt.expectedType {
				t.Fatalf("varStmt.Type wrong. expected=%s, got=%v", tt.expectedType, varStmt.Type)
			}
		}
	}
}

func TestShortVarAndAssign(t *testing.T) {
	input := `
name := "Muhammad"
name = "RustX"
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got=%d", len(program.Statements))
	}

	shortStmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt 0 is not ShortVarStatement, got=%T", program.Statements[0])
	}
	if shortStmt.Name.Value != "name" {
		t.Fatalf("shortStmt name wrong: got %s", shortStmt.Name.Value)
	}

	assignStmt, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt 1 is not AssignStatement, got=%T", program.Statements[1])
	}
	ident, ok := assignStmt.Left.(*ast.Identifier)
	if !ok || ident.Value != "name" {
		t.Fatalf("assignStmt target wrong")
	}
}

func TestFunctionLiteral(t *testing.T) {
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
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement, got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("not FunctionLiteral, got=%T", exprStmt.Expression)
	}

	if fn.Name == nil || fn.Name.Value != "add" {
		t.Fatalf("fn name wrong, got=%v", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn params count wrong, got=%d", len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Name != "int" {
		t.Fatalf("fn return type wrong, got=%v", fn.ReturnType)
	}
}

func TestForLoop(t *testing.T) {
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
		t.Fatalf("expected 2 statements, got=%d", len(program.Statements))
	}

	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got=%T", program.Statements[0])
	}
	for1, ok := stmt1.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expected ForExpression, got=%T", stmt1.Expression)
	}
	if for1.Init == nil || for1.Condition == nil || for1.Post == nil {
		t.Fatalf("expected 3-part for loop components to be non-nil")
	}

	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got=%T", program.Statements[1])
	}
	for2, ok := stmt2.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expected ForExpression, got=%T", stmt2.Expression)
	}
	if for2.Init != nil || for2.Condition == nil || for2.Post != nil {
		t.Fatalf("expected conditional for loop")
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true == true", "(true == true)"},
		{"a && b || c", "((a && b) || c)"},
		{"a || b && c", "(a || (b && c))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, err := range errs {
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
