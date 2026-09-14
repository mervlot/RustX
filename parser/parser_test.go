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
name := "RustX";
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
		{"name"},
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
			t.Errorf("stmt.Name.Value not '%s'. got='%s'", name, stmt.Name.Value)
			return false
		}
	case *ast.ExpressionStatement:
		assign, ok := stmt.Expression.(*ast.AssignmentExpression)
		if !ok {
			t.Errorf("stmt.Expression not *ast.AssignmentExpression. got=%T", stmt.Expression)
			return false
		}
		if assign.Name.Value != name {
			t.Errorf("assign.Name.Value not '%s'. got='%s'", name, assign.Name.Value)
			return false
		}
	default:
		t.Errorf("s not *ast.VarStatement or *ast.ExpressionStatement. got=%T", s)
		return false
	}
	return true
}

func TestForAndWhileStatements(t *testing.T) {
	input := `
while x < 10 {
	x = x + 1;
}

for i := 0; i < 5; i++ {
	println(i);
}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
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
