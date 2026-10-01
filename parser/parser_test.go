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
	var name string = "RustX";
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
		{"name", "string"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.VarStatement. got=%T", stmt)
		}

		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}

		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
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

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("forStmt init/condition/post missing")
	}

	whileStmt, ok := program.Statements[1].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[1] is not *ast.WhileStatement. got=%T", program.Statements[1])
	}
	if whileStmt.Condition == nil || whileStmt.Body == nil {
		t.Fatalf("whileStmt condition/body missing")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType != "int" {
		t.Fatalf("function return type wrong. got=%s", fn.ReturnType)
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
