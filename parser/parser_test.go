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
		expectedValue      interface{}
	}{
		{"var x int = 5;", "x", "int", 5},
		{"var y = true;", "y", "", true},
		{"var foobar float = 3.14;", "foobar", "float", 3.14},
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

		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestShorthandAndAssignStatements(t *testing.T) {
	input := `
	x := 10;
	x = 20;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	stmt1, ok := program.Statements[0].(*ast.ShorthandAssignStatement)
	if !ok {
		t.Fatalf("stmt1 is not *ast.ShorthandAssignStatement. got=%T", program.Statements[0])
	}
	if stmt1.Name.Value != "x" {
		t.Errorf("stmt1.Name.Value not 'x'. got=%q", stmt1.Name.Value)
	}

	stmt2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt2 is not *ast.AssignStatement. got=%T", program.Statements[1])
	}
	if stmt2.Name.Value != "x" {
		t.Errorf("stmt2.Name.Value not 'x'. got=%q", stmt2.Name.Value)
	}
}

func TestForLoopStatements(t *testing.T) {
	input := `
	for i := 0; i < 10; i++ {
		println(i);
	}
	for x < 5 {
		println(x);
	}
	for {
		println("inf");
	}
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// 1st loop
	loop1, ok := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("statement 1 is not ForExpression. got=%T", program.Statements[0])
	}
	if loop1.Init == nil {
		t.Errorf("loop1.Init is nil")
	}
	if loop1.Cond == nil {
		t.Errorf("loop1.Cond is nil")
	}
	if loop1.Post == nil {
		t.Errorf("loop1.Post is nil")
	}

	// 2nd loop
	loop2, ok := program.Statements[1].(*ast.ExpressionStatement).Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("statement 2 is not ForExpression. got=%T", program.Statements[1])
	}
	if loop2.Init != nil {
		t.Errorf("loop2.Init should be nil")
	}
	if loop2.Cond == nil {
		t.Errorf("loop2.Cond should not be nil")
	}

	// 3rd loop
	loop3, ok := program.Statements[2].(*ast.ExpressionStatement).Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("statement 3 is not ForExpression. got=%T", program.Statements[2])
	}
	if loop3.Init != nil || loop3.Cond != nil || loop3.Post != nil {
		t.Errorf("loop3 should be infinite loop, all parts nil")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err)
	}
	t.FailNow()
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

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		return testIdentifierOrString(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	f, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}

	if f.Value != value {
		t.Errorf("f.Value not %f. got=%f", value, f.Value)
		return false
	}

	return true
}

func testIdentifierOrString(t *testing.T, exp ast.Expression, value string) bool {
	switch e := exp.(type) {
	case *ast.Identifier:
		if e.Value != value {
			t.Errorf("ident.Value not %s. got=%s", value, e.Value)
			return false
		}
		return true
	case *ast.StringLiteral:
		if e.Value != value {
			t.Errorf("str.Value not %s. got=%s", value, e.Value)
			return false
		}
		return true
	}
	t.Errorf("exp not *ast.Identifier or *ast.StringLiteral. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	return true
}
