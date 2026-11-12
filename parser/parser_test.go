package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5
var y int = 10
var z string = "hello"
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
		{"x", ""},
		{"y", "int"},
		{"z", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got='%s'", tt.expectedType, varStmt.Type)
		}
	}
}

func TestAssignStatements(t *testing.T) {
	input := `
x := 42
x = 100
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	s1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok || !s1.IsDefine {
		t.Fatalf("stmt[0] expected define assign statement, got=%T, isDefine=%v", program.Statements[0], s1.IsDefine)
	}

	s2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok || s2.IsDefine {
		t.Fatalf("stmt[1] expected assign statement, got=%T, isDefine=%v", program.Statements[1], s2.IsDefine)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if function.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%q", function.Name.Value)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	if function.ReturnType != "int" {
		t.Fatalf("function return type wrong. want 'int', got=%q", function.ReturnType)
	}
}

func TestControlFlowParsing(t *testing.T) {
	input := `
if x < 10 {
	println("small")
} else if x == 10 {
	println("ten")
} else {
	println("large")
}

while count > 0 {
	count--
}

for i := 0; i < 5; i++ {
	println(i)
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a * b + c",
			"((a * b) + c)",
		},
		{
			"a == b && c != d || e < f",
			"(((a == b) && (c != d)) || (e < f))",
		},
		{
			"a[0] + 1",
			"((a[0]) + 1)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.Statements[0].(*ast.ExpressionStatement).Expression.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
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
