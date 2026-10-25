package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y = 10;
var foobar = 838383;
name := "Muhammad";
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
		{"name"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarOrAssignStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarOrAssignStatement(t *testing.T, s ast.Statement, name string) bool {
	switch stmt := s.(type) {
	case *ast.VarStatement:
		if stmt.Name.Value != name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", name, stmt.Name.Value)
			return false
		}
		return true
	case *ast.ExpressionStatement:
		if assignExp, ok := stmt.Expression.(*ast.AssignExpression); ok {
			if ident, ok := assignExp.Left.(*ast.Identifier); ok {
				if ident.Value != name {
					t.Errorf("ident.Value not '%s'. got=%s", name, ident.Value)
					return false
				}
				return true
			}
		}
		t.Errorf("stmt.Expression is not *ast.AssignExpression. got=%T", stmt.Expression)
		return false
	default:
		t.Errorf("s not *ast.VarStatement or ExpressionStatement (*ast.AssignExpression). got=%T", s)
		return false
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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
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

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type != "int" {
		t.Errorf("param 0 wrong. got=%s %s", fn.Parameters[0].Name.Value, fn.Parameters[0].Type)
	}

	if fn.ReturnType != "int" {
		t.Errorf("return type wrong. expected 'int', got=%s", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `for i := 0; i < 10; i++ { println(i); }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
	}

	if forExp.Init == nil || forExp.Condition == nil || forExp.Post == nil {
		t.Fatalf("for loop components missing")
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
