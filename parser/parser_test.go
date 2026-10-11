package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5;
var y float = 10.5;
var name string = "RustX";
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
		{"x", "int"},
		{"y", "float"},
		{"name", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.VarStatement. got=%T", stmt)
		}

		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}

		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShorthandAssignments(t *testing.T) {
	input := `
x := 10;
y = 20;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	assignStmt, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("statement[0] is not *ast.AssignStatement. got=%T", program.Statements[0])
	}

	if assignStmt.Name.Value != "x" {
		t.Errorf("assignStmt.Name.Value not 'x'. got=%s", assignStmt.Name.Value)
	}

	if assignStmt.Token.Literal != ":=" {
		t.Errorf("assignStmt.Token.Literal not ':='. got=%s", assignStmt.Token.Literal)
	}
}

func TestFunctionDeclaration(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expression is not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("fn.Name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestControlFlow(t *testing.T) {
	input := `
if x < 10 {
	y := 1;
} else {
	y := 2;
}

while x < 20 {
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
		t.Fatalf("program.Statements length not 3. got=%d", len(program.Statements))
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
