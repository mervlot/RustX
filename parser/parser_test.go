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
var foobar = 838383;
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
		{"foobar", ""},
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
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("statement is not ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fn.Name.Value != "add" {
		t.Fatalf("function name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. expected 2, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("parameter 0 wrong")
	}

	if fn.ReturnType != "int" {
		t.Errorf("return type wrong. expected 'int', got=%s", fn.ReturnType)
	}
}

func TestForStatement(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i); }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("statement is not ast.ForStatement. got=%T", program.Statements[0])
	}

	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("for loop structure missing parts: init=%v, cond=%v, post=%v",
			forStmt.Init, forStmt.Condition, forStmt.Post)
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
