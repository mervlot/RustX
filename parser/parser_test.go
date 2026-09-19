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
name := "RustX";
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
		{"name"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		switch s := stmt.(type) {
		case *ast.VarStatement:
			if s.Name.Value != tt.expectedIdentifier {
				t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, s.Name.Value)
			}
		case *ast.ExpressionStatement:
			infix, ok := s.Expression.(*ast.InfixExpression)
			if !ok || infix.Operator != ":=" {
				t.Fatalf("stmt is not := expression. got=%T", s.Expression)
			}
			ident, ok := infix.Left.(*ast.Identifier)
			if !ok || ident.Value != tt.expectedIdentifier {
				t.Errorf("ident.Value not '%s'. got=%s", tt.expectedIdentifier, ident.Value)
			}
		default:
			t.Fatalf("stmt not *ast.VarStatement or shorthand :=. got=%T", stmt)
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

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Fatalf("function name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters length wrong. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Name != "int" {
		t.Fatalf("return type not 'int'. got=%v", fn.ReturnType)
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

	stmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt is not ForStatement. got=%T", program.Statements[0])
	}

	if stmt.Init == nil || stmt.Condition == nil || stmt.Post == nil {
		t.Fatalf("for statement header components missing")
	}
}

func TestWhileStatementParsing(t *testing.T) {
	input := `while count < 5 { count++ }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt is not WhileStatement. got=%T", program.Statements[0])
	}

	if stmt.Condition == nil || stmt.Body == nil {
		t.Fatalf("while statement condition or body is nil")
	}
}

func TestIfElseParsing(t *testing.T) {
	input := `if x > 10 { println("high") } else if x > 5 { println("mid") } else { println("low") }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ExpressionStatement. got=%T", program.Statements[0])
	}

	ifExpr, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expr is not IfExpression. got=%T", stmt.Expression)
	}

	if ifExpr.Alternative == nil {
		t.Fatalf("ifExpr.Alternative is nil")
	}
}

func TestPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a && b || c",
			"((a && b) || c)",
		},
		{
			"a || b && c",
			"(a || (b && c))",
		},
		{
			"x == 5 && y > 10",
			"((x == 5) && (y > 10))",
		},
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
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
