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
var foobar int
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
		{"y", ""},
		{"foobar", "int"},
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
		if varStmt.TypeAnnotation != tt.expectedType {
			t.Errorf("varStmt.TypeAnnotation not '%s'. got=%s", tt.expectedType, varStmt.TypeAnnotation)
		}
	}
}

func TestShortVarAndAssignStatements(t *testing.T) {
	input := `
x := 5
x = 10
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	shortStmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}
	if shortStmt.Name.Value != "x" {
		t.Errorf("shortStmt.Name.Value not 'x'. got=%s", shortStmt.Name.Value)
	}

	assignStmt, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.AssignStatement. got=%T", program.Statements[1])
	}
	if assignStmt.Name.Value != "x" {
		t.Errorf("assignStmt.Name.Value not 'x'. got=%s", assignStmt.Name.Value)
	}
}

func TestFunctionStatement(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionStatement. got=%T", program.Statements[0])
	}

	if fn.Name.Value != "add" {
		t.Errorf("fn.Name.Value not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length not 2. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForLoopStatement(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}
for x < 5 {
	x++
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	for3Part, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if for3Part.Init == nil || for3Part.Condition == nil || for3Part.Post == nil {
		t.Errorf("for3Part init, condition or post is nil")
	}

	forSingle, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.ForStatement. got=%T", program.Statements[1])
	}
	if forSingle.Init != nil || forSingle.Condition == nil || forSingle.Post != nil {
		t.Errorf("forSingle header mismatch: init=%v, cond=%v, post=%v", forSingle.Init, forSingle.Condition, forSingle.Post)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}
