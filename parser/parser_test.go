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
	var foobar = true;
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
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, typeName string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if typeName != "" {
		if varStmt.Type == nil {
			t.Errorf("varStmt.Type is nil. expected=%s", typeName)
			return false
		}
		if varStmt.Type.Value != typeName {
			t.Errorf("varStmt.Type.Value not '%s'. got=%s", typeName, varStmt.Type.Value)
			return false
		}
	} else {
		if varStmt.Type != nil {
			t.Errorf("varStmt.Type is not nil. got=%s", varStmt.Type.Value)
			return false
		}
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a && b || c",
			"((a && b) || c)",
		},
		{
			"a || b && c",
			"(a || (b && c))",
		},
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a * b + c",
			"((a * b) + c)",
		},
		{
			"a++ + b",
			"((a++) + b)",
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

func TestForLoopParsing(t *testing.T) {
	input := `
	for i := 0; i < 10; i++ {
		println(i);
	}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected expression statement, got=%T", program.Statements[0])
	}

	fe, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expected for expression, got=%T", stmt.Expression)
	}

	if fe.Init == nil {
		t.Errorf("expected ForExpression.Init to be populated")
	}
	if fe.Condition == nil {
		t.Errorf("expected ForExpression.Condition to be populated")
	}
	if fe.Post == nil {
		t.Errorf("expected ForExpression.Post to be populated")
	}
}

func TestSyntaxErrorPosition(t *testing.T) {
	input := `
	var x = ;
	`
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Fatalf("expected parser errors, got none")
	}

	err := p.Errors()[0]
	expected := "SyntaxError:\nno prefix parse function for ; found\nat line 2 column 10"
	if err.Error() != expected {
		t.Errorf("expected error string:\n%q\ngot:\n%q", expected, err.Error())
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
