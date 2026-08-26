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
z := 15;
a = 20;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	expectedNames := []string{"x", "y", "z", "a"}
	for i, name := range expectedNames {
		stmt := program.Statements[i]
		switch s := stmt.(type) {
		case *ast.VarStatement:
			if s.Name.Value != name {
				t.Errorf("stmt.Name.Value not '%s'. got=%s", name, s.Name.Value)
			}
		case *ast.ShortVarStatement:
			if s.Name.Value != name {
				t.Errorf("stmt.Name.Value not '%s'. got=%s", name, s.Name.Value)
			}
		case *ast.AssignmentStatement:
			if s.Name.Value != name {
				t.Errorf("stmt.Name.Value not '%s'. got=%s", name, s.Name.Value)
			}
		default:
			t.Fatalf("unexpected statement type for index %d: %T", i, stmt)
		}
	}
}

func TestFunctionDeclarationAndCall(t *testing.T) {
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

	fnStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fnLit, ok := fnStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("fnStmt.Expression is not ast.FunctionLiteral. got=%T", fnStmt.Expression)
	}

	if fnLit.Name.Value != "add" {
		t.Errorf("fnLit.Name not 'add'. got=%s", fnLit.Name.Value)
	}

	if len(fnLit.Parameters) != 2 {
		t.Fatalf("fnLit.Parameters length not 2. got=%d", len(fnLit.Parameters))
	}
}

func TestControlFlowParsing(t *testing.T) {
	input := `
if x > 5 {
	println("greater");
} else {
	println("smaller");
}

while x < 10 {
	x++;
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
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %s", msg)
	}
	t.FailNow()
}
