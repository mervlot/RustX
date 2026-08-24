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

	stmt1, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt1 not *ast.VarStatement. got=%T", program.Statements[0])
	}
	if stmt1.Name.Value != "x" || stmt1.Type != "int" {
		t.Errorf("stmt1 incorrect. name=%s type=%s", stmt1.Name.Value, stmt1.Type)
	}

	stmt2, ok := program.Statements[1].(*ast.VarStatement)
	if !ok {
		t.Fatalf("stmt2 not *ast.VarStatement. got=%T", program.Statements[1])
	}
	if stmt2.Name.Value != "y" {
		t.Errorf("stmt2 name incorrect. name=%s", stmt2.Name.Value)
	}

	stmt3, ok := program.Statements[2].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt3 not *ast.ShortVarStatement. got=%T", program.Statements[2])
	}
	if stmt3.Name.Value != "name" {
		t.Errorf("stmt3 name incorrect. name=%s", stmt3.Name.Value)
	}
}

func TestFunctionStatementAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}
add(5, 10);
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if fn.Name.Value != "add" {
		t.Errorf("function name is not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. got=%d", len(fn.Parameters))
	}
}

func TestControlFlowParsing(t *testing.T) {
	input := `
if x < 10 {
	println(x);
} else {
	println(0);
}

while x > 0 {
	x--;
}

for i := 0; i < 5; i++ {
	println(i);
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements expected 3 statements, got=%d", len(program.Statements))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %q", msg.Error())
	}
	t.FailNow()
}
