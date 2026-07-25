package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y int = 10
var z float = 12.34;
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
		{"x", ""},
		{"y", "int"},
		{"z", "float"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, typ string) bool {
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Type != typ {
		t.Errorf("varStmt.Type not '%s'. got=%s", typ, varStmt.Type)
		return false
	}

	return true
}

func TestForStatements(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	x := 5
}
for x < 5 {
	x++
}
for {
	break
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// 1st is ForStatement (3 part)
	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt1 is not ExpressionStatement. got=%T", program.Statements[0])
	}
	forExp1, ok := stmt1.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt1.Expression is not ForExpression. got=%T", stmt1.Expression)
	}
	if forExp1.Initializer == nil || forExp1.Condition == nil || forExp1.PostCycle == nil {
		t.Fatalf("forExp1 has nil fields in 3-part loop")
	}

	// 2nd is ForStatement (condition)
	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt2 is not ExpressionStatement. got=%T", program.Statements[1])
	}
	forExp2, ok := stmt2.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt2.Expression is not ForExpression. got=%T", stmt2.Expression)
	}
	if forExp2.Initializer != nil || forExp2.Condition == nil || forExp2.PostCycle != nil {
		t.Fatalf("forExp2 should only have condition")
	}

	// 3rd is Infinite loop
	stmt3, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt3 is not ExpressionStatement. got=%T", program.Statements[2])
	}
	forExp3, ok := stmt3.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt3.Expression is not ForExpression. got=%T", stmt3.Expression)
	}
	if forExp3.Initializer != nil || forExp3.Condition != nil || forExp3.PostCycle != nil {
		t.Fatalf("forExp3 should be infinite loop")
	}
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
