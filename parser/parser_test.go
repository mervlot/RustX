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
var foobar = 838383;
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
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if stmt.TokenLiteral() != "var" {
			t.Fatalf("stmt.TokenLiteral not 'var'. got=%q", stmt.TokenLiteral())
		}
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
	}
}

func TestShorthandVarStatement(t *testing.T) {
	input := `name := "Muhammad";`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ShorthandVarStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ShorthandVarStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "name" {
		t.Fatalf("stmt.Name.Value not 'name'. got=%s", stmt.Name.Value)
	}
}

func TestFunctionDeclaration(t *testing.T) {
	input := `func add(a int, b int) int { return a + b }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Fatalf("stmt.Name.Value not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("stmt.Parameters length not 2. got=%d", len(stmt.Parameters))
	}
}

func TestForLoop(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i) }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ForStatement. got=%T", program.Statements[0])
	}
}

func TestWhileLoop(t *testing.T) {
	input := `while count > 0 { count-- }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt not *ast.WhileStatement. got=%T", program.Statements[0])
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a + b * c", "(a + (b * c))"},
		{"a && b || c", "((a && b) || c)"},
		{"a == b && c != d", "((a == b) && (c != d))"},
		{"a + b[1] * c", "(a + ((b[1]) * c))"},
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
