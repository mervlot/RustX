package parser

import (
	"rustx/ast"
	"rustx/lexer"
	"testing"
)

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y int = 10;
var z: string = "hello";
name := "Muhammad";
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"z"},
		{"name"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not *ast.VarStatement. got=%T", i, stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
	}
}

func TestFunctionParsing(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain 1 statement. got=%d", len(program.Statements))
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
		t.Fatalf("function name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters wrong count. got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "a" || fn.Parameters[0].Type.Value != "int" {
		t.Fatalf("param 0 wrong. got=%+v", fn.Parameters[0])
	}
	if fn.Parameters[1].Name.Value != "b" || fn.Parameters[1].Type.Value != "int" {
		t.Fatalf("param 1 wrong. got=%+v", fn.Parameters[1])
	}

	if fn.ReturnType == nil || fn.ReturnType.Value != "int" {
		t.Fatalf("return type wrong. got=%+v", fn.ReturnType)
	}
}

func TestControlFlow(t *testing.T) {
	input := `
if x < 10 {
	y = 1
} else {
	y = 2
}

while count > 0 {
	count--
}

for i := 0; i < 10; i++ {
	println(i)
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got=%d", len(program.Statements))
	}

	// 1. If statement
	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 0 is not ExpressionStatement, got=%T", program.Statements[0])
	}
	ifExp, ok := stmt1.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt 0 expr is not IfExpression, got=%T", stmt1.Expression)
	}
	if ifExp.Alternative == nil {
		t.Fatalf("expected else branch in IfExpression")
	}

	// 2. While statement
	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 1 is not ExpressionStatement, got=%T", program.Statements[1])
	}
	_, ok = stmt2.Expression.(*ast.WhileExpression)
	if !ok {
		t.Fatalf("stmt 1 expr is not WhileExpression, got=%T", stmt2.Expression)
	}

	// 3. For statement
	stmt3, ok := program.Statements[2].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 2 is not ExpressionStatement, got=%T", program.Statements[2])
	}
	forExp, ok := stmt3.Expression.(*ast.ForExpression)
	if !ok {
		t.Fatalf("stmt 2 expr is not ForExpression, got=%T", stmt3.Expression)
	}
	if forExp.Init == nil || forExp.Condition == nil || forExp.Post == nil {
		t.Fatalf("forExp components missing: init=%v, cond=%v, post=%v", forExp.Init, forExp.Condition, forExp.Post)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b * c", "(a + (b * c))"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true && false || true", "((true && false) || true)"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"a * [1, 2, 3, 4][b * c] * d", "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg.Error())
	}
	t.FailNow()
}
