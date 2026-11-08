package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

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

func TestVarStatements(t *testing.T) {
	input := `
var x = 5
var y int = 10
var foobar = true
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
		{"foobar", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Fatalf("varStmt.Type not '%s'. got='%s'", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShortVarStatement(t *testing.T) {
	input := `name := "Muhammad"`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}
	if stmt.Name.Value != "name" {
		t.Fatalf("name wrong: got=%s", stmt.Name.Value)
	}
	strLit, ok := stmt.Value.(*ast.StringLiteral)
	if !ok || strLit.Value != "Muhammad" {
		t.Fatalf("value wrong: got=%v", stmt.Value)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}

	if function.Name == nil || function.Name.Value != "add" {
		t.Fatalf("function name not 'add'")
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	if function.ReturnType != "int" {
		t.Fatalf("function return type wrong. want 'int', got=%s", function.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i) }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := exprStmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expr is not ast.ForExpression. got=%T", exprStmt.Expression)
	}

	if forExp.Init == nil || forExp.Condition == nil || forExp.Post == nil {
		t.Fatalf("expected 3-part for loop, got nil parts: init=%v, cond=%v, post=%v",
			forExp.Init, forExp.Condition, forExp.Post)
	}
}

func TestWhileLoopParsing(t *testing.T) {
	input := `while x < 10 { x++ }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("not ExpressionStatement. got=%T", program.Statements[0])
	}

	whileExp, ok := exprStmt.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("not WhileExpression. got=%T", exprStmt.Expression)
	}

	if whileExp.Condition == nil || whileExp.Body == nil {
		t.Fatalf("while expression incomplete")
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b * c",
			"(a + (b * c))\n",
		},
		{
			"a && b || c",
			"((a && b) || c)\n",
		},
		{
			"a || b && c",
			"(a || (b && c))\n",
		},
		{
			"5 < 10 == true",
			"((5 < 10) == true)\n",
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
