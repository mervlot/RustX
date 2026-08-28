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
var foobar float = 8.3;
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
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
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

	return true
}

func TestAssignStatements(t *testing.T) {
	input := `
x := 5
y = 10
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedName string
		operator     string
	}{
		{"x", ":="},
		{"y", "="},
	}

	for i, tt := range tests {
		stmt, ok := program.Statements[i].(*ast.AssignStatement)
		if !ok {
			t.Fatalf("stmt[%d] not *ast.AssignStatement. got=%T", i, program.Statements[i])
		}
		if stmt.Name.Value != tt.expectedName {
			t.Errorf("stmt[%d].Name.Value not %s. got=%s", i, tt.expectedName, stmt.Name.Value)
		}
		if stmt.Token.Literal != tt.operator {
			t.Errorf("stmt[%d].Token.Literal not %s. got=%s", i, tt.operator, stmt.Token.Literal)
		}
	}
}

func TestForAndWhileStatements(t *testing.T) {
	input := `
while x < 10 {
	x++
}

for i := 0; i < 10; i++ {
	println(i)
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	_, okWhile := program.Statements[0].(*ast.WhileStatement)
	if !okWhile {
		t.Fatalf("stmt[0] not *ast.WhileStatement. got=%T", program.Statements[0])
	}

	_, okFor := program.Statements[1].(*ast.ForStatement)
	if !okFor {
		t.Fatalf("stmt[1] not *ast.ForStatement. got=%T", program.Statements[1])
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b; }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fn.Name.Value != "add" {
		t.Fatalf("function name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count not 2. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Value != "int" {
		t.Fatalf("return type not 'int'")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg.Error())
	}
	t.FailNow()
}
