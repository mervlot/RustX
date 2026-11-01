package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
	var x int = 5;
	var y = 10;
	var foobar = true;
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
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got='%s'",
				tt.expectedIdentifier, varStmt.Name.Value)
		}
	}
}

func TestAssignAndShorthand(t *testing.T) {
	input := `
	x := 42
	x = 100
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	s1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok || !s1.IsDefine {
		t.Fatalf("expected AssignStatement with IsDefine=true")
	}

	s2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok || s2.IsDefine {
		t.Fatalf("expected AssignStatement with IsDefine=false")
	}
}

func TestFunctionAndControlFlow(t *testing.T) {
	input := `
	func add(a int, b int) int {
		return a + b
	}

	if x < 10 {
		println("low")
	} else {
		println("high")
	}

	while x > 0 {
		x--
	}

	for i := 0; i < 10; i++ {
		println(i)
	}
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("expected 4 statements, got=%d", len(program.Statements))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}
