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
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShorthandAndAssignStatements(t *testing.T) {
	input := `
x := 5;
x = 10;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	for i, stmt := range program.Statements {
		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Fatalf("stmt[%d] not *ast.AssignStatement. got=%T", i, stmt)
		}
		if assignStmt.Name.Value != "x" {
			t.Errorf("assignStmt.Name.Value not 'x'. got=%s", assignStmt.Name.Value)
		}
	}
}

func TestControlFlowStatements(t *testing.T) {
	input := `
if x < y {
	return true;
} else {
	return false;
}

while x < 10 {
	x++;
}

for i := 0; i < 10; i++ {
	println(i);
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

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}

add(5, 10);
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fnStmt, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt[0] not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}
	if fnStmt.Name.Value != "add" {
		t.Errorf("fnStmt.Name.Value not 'add'. got=%s", fnStmt.Name.Value)
	}
	if fnStmt.ReturnType != "int" {
		t.Errorf("fnStmt.ReturnType not 'int'. got=%s", fnStmt.ReturnType)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
