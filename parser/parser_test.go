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
z := 15;
a = 20;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	// 1. var x int = 5
	varStmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt[0] is not *ast.VarStatement. got=%T", program.Statements[0])
	}
	if varStmt.Name.Value != "x" {
		t.Errorf("varStmt.Name.Value not 'x'. got=%q", varStmt.Name.Value)
	}
	if varStmt.Type == nil || varStmt.Type.Value != "int" {
		t.Errorf("varStmt.Type.Value not 'int'")
	}

	// 2. var y = 10
	varStmt2, ok := program.Statements[1].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt[1] is not *ast.VarStatement. got=%T", program.Statements[1])
	}
	if varStmt2.Name.Value != "y" {
		t.Errorf("varStmt2.Name.Value not 'y'. got=%q", varStmt2.Name.Value)
	}

	// 3. z := 15
	shortStmt, ok := program.Statements[2].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt[2] is not *ast.ShortVarStatement. got=%T", program.Statements[2])
	}
	if shortStmt.Name.Value != "z" {
		t.Errorf("shortStmt.Name.Value not 'z'. got=%q", shortStmt.Name.Value)
	}

	// 4. a = 20
	assignStmt, ok := program.Statements[3].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt[3] is not *ast.AssignmentStatement. got=%T", program.Statements[3])
	}
	if assignStmt.Name.Value != "a" {
		t.Errorf("assignStmt.Name.Value not 'a'. got=%q", assignStmt.Name.Value)
	}
}

func TestFunctionAndControlFlow(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}

if x > 0 {
	println("positive")
} else {
	println("non-positive")
}

for i := 0; i < 10; i++ {
	println(i)
}

while x < 20 {
	x++
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
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
