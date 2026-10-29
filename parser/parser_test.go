package parser

import (
	"testing"

	"rustx/ast"
	"rustx/lexer"
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
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
	}
}

func TestFunctionAndControlFlow(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}

for i := 0; i < 10; i++ {
	if i == 5 {
		println(i)
	}
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
		t.Errorf("parser error: %s", err.Error())
	}
	t.FailNow()
}
