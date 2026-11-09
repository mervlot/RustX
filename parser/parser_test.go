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
		hasType             bool
		expectedType       string
	}{
		{"x", false, ""},
		{"y", true, "int"},
		{"foobar", false, ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("program.Statements[%d] is not *ast.VarStatement. got=%T", i, stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if tt.hasType {
			if varStmt.Type == nil || varStmt.Type.Value != tt.expectedType {
				t.Fatalf("varStmt.Type expected %s, got %v", tt.expectedType, varStmt.Type)
			}
		}
	}
}

func TestAssignAndShorthand(t *testing.T) {
	input := `
x := 15;
x = 20;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got=%d", len(program.Statements))
	}

	s1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok || !s1.IsDefine || s1.Name.Value != "x" {
		t.Fatalf("expected define assign for x, got=%v", program.Statements[0])
	}

	s2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok || s2.IsDefine || s2.Name.Value != "x" {
		t.Fatalf("expected regular assign for x, got=%v", program.Statements[1])
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expected FunctionLiteral, got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Fatalf("expected function name 'add', got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("expected 2 parameters, got=%d", len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Value != "int" {
		t.Fatalf("expected return type int, got=%v", fn.ReturnType)
	}
}

func TestForLoopParsing(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
    println(i)
}
for x < 5 {
    x++
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got=%T", program.Statements[0])
	}

	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("expected init, cond, and post to be non-nil")
	}

	forStmt2, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got=%T", program.Statements[1])
	}
	if forStmt2.Init != nil || forStmt2.Condition == nil {
		t.Fatalf("expected cond only for second loop")
	}
}

func TestControlFlowParsing(t *testing.T) {
	input := `
if x > 5 {
    println("greater")
} else if x == 5 {
    println("equal")
} else {
    println("less")
}

while y > 0 {
    y--
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 2 statements, got=%d", len(program.Statements))
	}

	ifStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement for if, got=%T", program.Statements[0])
	}
	ifExp, ok := ifStmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expected IfExpression, got=%T", ifStmt.Expression)
	}
	if ifExp.Alternative == nil {
		t.Fatalf("expected alternative in ifExp")
	}

	whileStmt, ok := program.Statements[1].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("expected WhileStatement, got=%T", program.Statements[1])
	}
	if whileStmt.Condition == nil || whileStmt.Body == nil {
		t.Fatalf("expected condition and body in whileStmt")
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %s", msg.Error())
	}
	t.FailNow()
}
