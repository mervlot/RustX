package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5
var y = 10
name := "RustX"
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
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	switch stmt := s.(type) {
	case *ast.VarStatement:
		if stmt.Name.Value != name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
			return false
		}
	case *ast.ExpressionStatement:
		if assign, ok := stmt.Expression.(*ast.AssignmentExpression); ok {
			if assign.Name.Value != name {
				t.Errorf("assign.Name.Value not '%s'. got=%s", name, assign.Name.Value)
				return false
			}
		} else {
			t.Errorf("s not *ast.VarStatement or *ast.AssignmentExpression. got=%T", s)
			return false
		}
	default:
		t.Errorf("s not *ast.VarStatement or *ast.ExpressionStatement. got=%T", s)
		return false
	}
	return true
}

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
add(5, 10)
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}
}

func TestForAndWhileStatements(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}
while x > 0 {
	x--
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
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, err := range errs {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
