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
var z = true
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
		{"z", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string, expectedType string) bool {
	if s.GetToken().Literal != "var" {
		t.Errorf("s.GetToken().Literal not 'var'. got=%q", s.GetToken().Literal)
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

	if varStmt.Type != expectedType {
		t.Errorf("varStmt.Type not '%s'. got=%s", expectedType, varStmt.Type)
		return false
	}

	return true
}

func TestShorthandAndAssignExpression(t *testing.T) {
	input := `x := 42
y = 100`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("Expected 2 statements, got=%d", len(program.Statements))
	}

	// Test :=
	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement 1 not ExpressionStatement, got=%T", program.Statements[0])
	}
	assign1, ok := stmt1.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("Expression 1 not AssignmentExpression, got=%T", stmt1.Expression)
	}
	if assign1.Operator != ":=" {
		t.Errorf("Operator not ':=', got=%s", assign1.Operator)
	}

	// Test =
	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement 2 not ExpressionStatement, got=%T", program.Statements[1])
	}
	assign2, ok := stmt2.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("Expression 2 not AssignmentExpression, got=%T", stmt2.Expression)
	}
	if assign2.Operator != "=" {
		t.Errorf("Operator not '=', got=%s", assign2.Operator)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b * c", "(a + (b * c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"!(true == true)", "(!(true == true))"},
		{"a && b || c", "((a && b) || c)"},
		{"a || b && c", "(a || (b && c))"},
		{"a++", "(a++)"},
		{"b--", "(b--)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q for input=%q", tt.expected, actual, tt.input)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := "if x < y { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if exp.Condition.String() != "(x < y)" {
		t.Errorf("condition is not (x < y). got=%s", exp.Condition.String())
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(exp.Consequence.Statements))
	}

	if exp.Consequence.Statements[0].String() != "x" {
		t.Errorf("consequence statement is not 'x'. got=%s", exp.Consequence.Statements[0].String())
	}

	if exp.Alternative == nil {
		t.Fatalf("exp.Alternative is nil")
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement. got=%d", len(exp.Alternative.Statements))
	}

	if exp.Alternative.Statements[0].String() != "y" {
		t.Errorf("alternative statement is not 'y'. got=%s", exp.Alternative.Statements[0].String())
	}
}

func TestForStatements(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}
for x < 10 {
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
		t.Fatalf("Expected 3 statements, got=%d", len(program.Statements))
	}

	for1, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("First statement not *ast.ForStatement, got=%T", program.Statements[0])
	}
	if for1.Init == nil || for1.Condition == nil || for1.Post == nil {
		t.Errorf("First loop misses 3-part elements. got init=%v, cond=%v, post=%v", for1.Init, for1.Condition, for1.Post)
	}

	for2, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("Second statement not *ast.ForStatement, got=%T", program.Statements[1])
	}
	if for2.Init != nil || for2.Condition == nil || for2.Post != nil {
		t.Errorf("Second loop should only have condition. got init=%v, cond=%v, post=%v", for2.Init, for2.Condition, for2.Post)
	}

	for3, ok := program.Statements[2].(*ast.ForStatement)
	if !ok {
		t.Fatalf("Third statement not *ast.ForStatement, got=%T", program.Statements[2])
	}
	if for3.Init != nil || for3.Condition != nil || for3.Post != nil {
		t.Errorf("Third loop should have no header. got init=%v, cond=%v, post=%v", for3.Init, for3.Condition, for3.Post)
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `func add(x int, y int) int { return x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("function name wrong. expected 'add', got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("parameters wrong. expected 2, got=%d", len(fn.Parameters))
	}

	if fn.Parameters[0].Name.Value != "x" || fn.Parameters[0].Type != "int" {
		t.Errorf("param 0 wrong. got=%v", fn.Parameters[0])
	}

	if fn.ReturnType != "int" {
		t.Errorf("return type wrong. expected 'int', got=%s", fn.ReturnType)
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
