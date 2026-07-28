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
	var foobar = true;
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
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral not 'var'. got=%q", s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	return true
}

func TestShortAssignAndAssignment(t *testing.T) {
	input := `
	x := 5;
	x = 10;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	// First statement is directly a ShortAssignStatement since it is parsed as a Statement!
	shortAssign, ok := program.Statements[0].(*ast.ShortAssignStatement)
	if !ok {
		// Or if parsed as ExpressionStatement holding ShortAssignStatement
		if estmt, ok := program.Statements[0].(*ast.ExpressionStatement); ok {
			shortAssign, ok = estmt.Expression.(*ast.ShortAssignStatement)
			if !ok {
				t.Fatalf("stmt1 not *ast.ShortAssignStatement. got=%T", estmt.Expression)
			}
		} else {
			t.Fatalf("stmt1 not ShortAssignStatement or ExpressionStatement. got=%T", program.Statements[0])
		}
	}
	if shortAssign.Name.Value != "x" {
		t.Errorf("shortAssign.Name.Value not 'x'. got=%q", shortAssign.Name.Value)
	}

	// Second statement is directly an AssignmentStatement
	assign, ok := program.Statements[1].(*ast.AssignmentStatement)
	if !ok {
		if estmt, ok := program.Statements[1].(*ast.ExpressionStatement); ok {
			assign, ok = estmt.Expression.(*ast.AssignmentStatement)
			if !ok {
				t.Fatalf("stmt2 not *ast.AssignmentStatement. got=%T", estmt.Expression)
			}
		} else {
			t.Fatalf("stmt2 not AssignmentStatement or ExpressionStatement. got=%T", program.Statements[1])
		}
	}
	ident, ok := assign.Name.(*ast.Identifier)
	if !ok {
		t.Fatalf("assign.Name not *ast.Identifier. got=%T", assign.Name)
	}
	if ident.Value != "x" {
		t.Errorf("ident.Value not 'x'. got=%q", ident.Value)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `
	for i := 0; i < 10; i++ {
		println(i);
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
		t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExpr, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("expr not *ast.ForExpression. got=%T", stmt.Expression)
	}

	if forExpr.Init == nil {
		t.Fatalf("forExpr.Init is nil")
	}

	if forExpr.Condition == nil {
		t.Fatalf("forExpr.Condition is nil")
	}

	if forExpr.Post == nil {
		t.Fatalf("forExpr.Post is nil")
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
