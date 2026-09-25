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
name := "RustX";
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrorsTest(t, p)

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
		if !testVarOrColonAssign(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarOrColonAssign(t *testing.T, s ast.Statement, name string) bool {
	switch stmt := s.(type) {
	case *ast.VarStatement:
		if stmt.Name.Value != name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
			return false
		}
	case *ast.ColonAssignStatement:
		if stmt.Name.Value != name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
			return false
		}
	default:
		t.Errorf("s not *ast.VarStatement or *ast.ColonAssignStatement. got=%T", s)
		return false
	}
	return true
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrorsTest(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. expected 2, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type.Value != "int" {
		t.Fatalf("parameter 0 wrong. got=%s %s", fn.Parameters[0].Name.Value, fn.Parameters[0].Type.Value)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i); }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrorsTest(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if forStmt.Init == nil {
		t.Fatalf("forStmt.Init is nil")
	}
	if forStmt.Condition == nil {
		t.Fatalf("forStmt.Condition is nil")
	}
	if forStmt.Post == nil {
		t.Fatalf("forStmt.Post is nil")
	}
}

func checkParserErrorsTest(t *testing.T, p *Parser) {
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
