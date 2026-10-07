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
			t.Errorf("varStmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, varStmt.Name.Value)
		}
		if varStmt.ValueType != tt.expectedType {
			t.Errorf("varStmt.ValueType not '%s'. got=%s", tt.expectedType, varStmt.ValueType)
		}
	}
}

func TestAssignAndShorthand(t *testing.T) {
	input := `
x := 5;
x = 10;
arr[0] = 15;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		_, ok := stmt.(*ast.AssignStatement)
		if !ok {
			t.Fatalf("stmt not *ast.AssignStatement. got=%T", stmt)
		}
	}
}

func TestFunctionLiteralAndCall(t *testing.T) {
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

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	fn, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expr not *ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}

	if fn.Name != "add" {
		t.Errorf("fn.Name not 'add'. got=%s", fn.Name)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters length not 2. got=%d", len(fn.Parameters))
	}

	if fn.ParamTypes[0] != "int" || fn.ParamTypes[1] != "int" {
		t.Errorf("fn.ParamTypes wrong: %v", fn.ParamTypes)
	}

	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
	}
}

func TestForLoopAndControlFlow(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i);
}

while x < 5 {
	x++;
}

if x == 5 {
	println("five");
} else if x == 10 {
	println("ten");
} else {
	println("other");
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
		t.Errorf("parser error: %s", msg.Error())
	}
	t.FailNow()
}
