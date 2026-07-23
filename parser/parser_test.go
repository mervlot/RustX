package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedType       string
		expectedValue      string
	}{
		{"var x = 5;", "x", "", "5"},
		{"var y int = 10;", "y", "int", "10"},
		{"var name string = \"RustX\";", "name", "string", "RustX"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}

		varStmt := stmt.(*ast.VarStatement)
		if varStmt.Value.String() != tt.expectedValue {
			t.Errorf("varStmt.Value.String() not %q. got=%q", tt.expectedValue, varStmt.Value.String())
		}
	}
}

func TestForLoopExpression(t *testing.T) {
	input := `for i := 0; i < 10; i++ {
		println(i)
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	forExp, ok := stmt.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.ForExpression. got=%T", stmt.Expression)
	}

	// Verify Init
	initStmt, ok := forExp.Init.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("forExp.Init is not ast.ExpressionStatement. got=%T", forExp.Init)
	}
	colonAssign, ok := initStmt.Expression.(*ast.ColonAssignExpression)
	if !ok {
		t.Fatalf("initStmt.Expression is not ast.ColonAssignExpression. got=%T", initStmt.Expression)
	}
	if colonAssign.Left.String() != "i" || colonAssign.Value.String() != "0" {
		t.Errorf("colonAssign values wrong. got left=%q, right=%q", colonAssign.Left.String(), colonAssign.Value.String())
	}

	// Verify Condition
	cond, ok := forExp.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("forExp.Condition is not ast.InfixExpression. got=%T", forExp.Condition)
	}
	if cond.Left.String() != "i" || cond.Operator != "<" || cond.Right.String() != "10" {
		t.Errorf("condition wrong. got=%q", cond.String())
	}

	// Verify Post
	postStmt, ok := forExp.Post.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("forExp.Post is not ast.ExpressionStatement. got=%T", forExp.Post)
	}
	postfix, ok := postStmt.Expression.(*ast.PostfixExpression)
	if !ok {
		t.Fatalf("postStmt.Expression is not ast.PostfixExpression. got=%T", postStmt.Expression)
	}
	if postfix.Left.String() != "i" || postfix.Operator != "++" {
		t.Errorf("postfix wrong. got=%q", postfix.String())
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, typeAnn string) bool {
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

	if varStmt.Type != typeAnn {
		t.Errorf("varStmt.Type not '%s'. got=%s", typeAnn, varStmt.Type)
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err.Error())
	}
	t.FailNow()
}
