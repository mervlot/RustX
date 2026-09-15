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
var foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", "int"},
		{"y", ""},
		{"foobar", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not '%s'. got='%s'", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Fatalf("varStmt.Type not '%s'. got='%s'", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShorthandAndAssignments(t *testing.T) {
	input := `
x := 5;
y = 10;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	stmt1, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt1 not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	assign1, ok := stmt1.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("expr1 not *ast.AssignmentExpression. got=%T", stmt1.Expression)
	}
	if !assign1.Shorthand || assign1.Name.Value != "x" {
		t.Fatalf("assign1 invalid: %+v", assign1)
	}

	stmt2, ok := program.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt2 not *ast.ExpressionStatement. got=%T", program.Statements[1])
	}
	assign2, ok := stmt2.Expression.(*ast.AssignmentExpression)
	if !ok {
		t.Fatalf("expr2 not *ast.AssignmentExpression. got=%T", stmt2.Expression)
	}
	if assign2.Shorthand || assign2.Name.Value != "y" {
		t.Fatalf("assign2 invalid: %+v", assign2)
	}
}

func TestFunctionAndControlFlowParsing(t *testing.T) {
	input := `
func add(a int, b int) int {
	if a < b {
		return a + b;
	} else {
		return b;
	}
}

while x < 10 {
	x++;
}

for i := 0; i < 10; i++ {
	println(i);
}
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements expected 3 statements. got=%d", len(program.Statements))
	}

	// 1. Function statement/literal
	fnStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt 0 not ExpressionStatement. got=%T", program.Statements[0])
	}
	fn, ok := fnStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr 0 not FunctionLiteral. got=%T", fnStmt.Expression)
	}
	if fn.Name.Value != "add" || len(fn.Parameters) != 2 || fn.ReturnType != "int" {
		t.Fatalf("function add parsing failed: %+v", fn)
	}

	// 2. While statement
	whileStmt, ok := program.Statements[1].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt 1 not WhileStatement. got=%T", program.Statements[1])
	}
	if whileStmt.Condition == nil || len(whileStmt.Body.Statements) != 1 {
		t.Fatalf("while statement parsing failed")
	}

	// 3. For statement
	forStmt, ok := program.Statements[2].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt 2 not ForStatement. got=%T", program.Statements[2])
	}
	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Fatalf("for statement parsing failed")
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
