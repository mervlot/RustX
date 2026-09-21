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
z := 42;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
		{"z"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarOrAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarOrAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	switch stmt := s.(type) {
	case *ast.VarStatement:
		if stmt.Name.Value != name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
			return false
		}
		return true
	case *ast.ExpressionStatement:
		if ae, ok := stmt.Expression.(*ast.AssignmentExpression); ok {
			if ae.Name.Value != name {
				t.Errorf("ae.Name.Value not '%s'. got=%s", name, ae.Name.Value)
				return false
			}
			return true
		}
	}
	t.Errorf("s not *ast.VarStatement or assignment *ast.ExpressionStatement. got=%T", s)
	return false
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(x int, y int) int { return x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt is not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%s", stmt.Name.Value)
	}

	if len(stmt.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. got=%d", len(stmt.Parameters))
	}

	if stmt.ReturnType != "int" {
		t.Fatalf("function return type wrong. got=%s", stmt.ReturnType)
	}
}

func TestControlFlowStatements(t *testing.T) {
	input := `
if x < y {
	return x;
} else {
	return y;
}

while x > 0 {
	x = x - 1;
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
