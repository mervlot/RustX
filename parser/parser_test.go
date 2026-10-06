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
var foobar = 838383;
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
		{"foobar", ""},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		varStmt, ok := stmt.(*ast.VarStatement)
		if !ok {
			t.Fatalf("stmt not *ast.VarStatement. got=%T", stmt)
		}
		if varStmt.Name.Value != tt.expectedIdentifier {
			t.Fatalf("varStmt.Name.Value not %q. got=%q", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.Type != tt.expectedType {
			t.Fatalf("varStmt.Type not %q. got=%q", tt.expectedType, varStmt.Type)
		}
	}
}

func TestShorthandAndAssignments(t *testing.T) {
	input := `
x := 5;
y = x + 1;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	assign1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.AssignStatement. got=%T", program.Statements[0])
	}
	if assign1.Token.Literal != ":=" {
		t.Fatalf("assign1 token literal not :=. got=%q", assign1.Token.Literal)
	}

	assign2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.AssignStatement. got=%T", program.Statements[1])
	}
	if assign2.Token.Literal != "=" {
		t.Fatalf("assign2 token literal not =. got=%q", assign2.Token.Literal)
	}
}

func TestFunctionDeclarationAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}

result := add(5, 10);
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt[0] not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}
	if fn.Name.Value != "add" {
		t.Fatalf("fn.Name.Value not 'add'. got=%q", fn.Name.Value)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}
	if fn.ReturnType != "int" {
		t.Fatalf("fn.ReturnType not 'int'. got=%q", fn.ReturnType)
	}
}

func TestControlFlow(t *testing.T) {
	input := `
if x < y {
	println(x);
} else {
	println(y);
}

while i < 10 {
	i++;
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
		t.Fatalf("program.Statements count wrong. got=%d", len(program.Statements))
	}

	// 1. If statement
	stmt0, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	_, ok = stmt0.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt[0].Expression not *ast.IfExpression. got=%T", stmt0.Expression)
	}

	// 2. While statement
	_, ok = program.Statements[1].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.WhileStatement. got=%T", program.Statements[1])
	}

	// 3. For statement
	_, ok = program.Statements[2].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[2] not *ast.ForStatement. got=%T", program.Statements[2])
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
