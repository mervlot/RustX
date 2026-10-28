package parser

import (
	"testing"

	"rustx/ast"
	"rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5
var y int = 10
var name string = "RustX"
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
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
	}
}

func TestShorthandAssignments(t *testing.T) {
	input := `
x := 5
y := "Hello"
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		assignStmt, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Fatalf("stmt not *ast.AssignStatement. got=%T", stmt)
		}
		if !assignStmt.IsDefine {
			t.Errorf("expected shorthand assign IsDefine to be true")
		}
	}
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

	fnStmt, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fnStmt.Name.Value != "add" {
		t.Errorf("fnStmt.Name.Value expected 'add', got=%s", fnStmt.Name.Value)
	}

	if len(fnStmt.Parameters) != 2 {
		t.Fatalf("fnStmt.Parameters length wrong. expected 2, got=%d", len(fnStmt.Parameters))
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
