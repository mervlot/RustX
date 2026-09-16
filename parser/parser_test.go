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
z := 15
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
		{"z"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarOrShortAssign(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarOrShortAssign(t *testing.T, s ast.Statement, name string) bool {
	switch stmt := s.(type) {
	case *ast.VarStatement:
		if stmt.Name.Value != name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
			return false
		}
		return true
	case *ast.ExpressionStatement:
		sae, ok := stmt.Expression.(*ast.ShortAssignExpression)
		if !ok {
			t.Errorf("stmt.Expression not *ast.ShortAssignExpression. got=%T", stmt.Expression)
			return false
		}
		if sae.Name.Value != name {
			t.Errorf("sae.Name.Value not '%s'. got=%s", name, sae.Name.Value)
			return false
		}
		return true
	default:
		t.Errorf("s not *ast.VarStatement or *ast.ExpressionStatement. got=%T", s)
		return false
	}
}

func TestFunctionLiteralAndControlFlow(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}

if x > 0 {
	println(x)
} else {
	println(0)
}

for i := 0; i < 10; i++ {
	println(i)
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
	for _, err := range errors {
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
