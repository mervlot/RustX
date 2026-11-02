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
var z = true
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		hasType            bool
	}{
		{"x", false},
		{"y", true},
		{"z", false},
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
		if tt.hasType && varStmt.Type == nil {
			t.Errorf("expected type for statement %d", i)
		}
	}
}

func TestAssignStatements(t *testing.T) {
	input := `
x := 5
x = 10
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	s0 := program.Statements[0].(*ast.AssignStatement)
	if !s0.IsDefine {
		t.Errorf("expected s0.IsDefine to be true")
	}

	s1 := program.Statements[1].(*ast.AssignStatement)
	if s1.IsDefine {
		t.Errorf("expected s1.IsDefine to be false")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(x int, y int) int { return x + y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if fn.Name != "add" {
		t.Fatalf("fn.Name wrong. expected='add', got=%q", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Name != "int" {
		t.Fatalf("function return type wrong. got=%v", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
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
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T",
			program.Statements[0])
	}

	if forStmt.Init == nil {
		t.Errorf("forStmt.Init is nil")
	}
	if forStmt.Condition == nil {
		t.Errorf("forStmt.Condition is nil")
	}
	if forStmt.Post == nil {
		t.Errorf("forStmt.Post is nil")
	}
}

func TestWhileLoopParsing(t *testing.T) {
	input := `
while x < 10 {
	x++
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"true && false || true",
			"((true && false) || true)",
		},
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
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
