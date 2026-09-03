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
var name string = "RustX";
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
		{"x", "int"},
		{"y", ""},
		{"name", "string"},
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
		if varStmt.Type != tt.expectedType {
			t.Errorf("varStmt.Type not '%s'. got=%s", tt.expectedType, varStmt.Type)
		}
	}
}

func TestForAndWhileStatements(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i);
}

while x < 5 {
	x = x + 1;
}
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	_, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt[0] not *ast.ForStatement. got=%T", program.Statements[0])
	}

	_, ok = program.Statements[1].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt[1] not *ast.WhileStatement. got=%T", program.Statements[1])
	}
}

func TestFunctionDeclaration(t *testing.T) {
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
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	fn, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fn.Name != "add" {
		t.Errorf("fn.Name not 'add'. got=%s", fn.Name)
	}
	if len(fn.Parameters) != 2 {
		t.Fatalf("fn.Parameters count wrong. got=%d", len(fn.Parameters))
	}
	if fn.ReturnType != "int" {
		t.Errorf("fn.ReturnType not 'int'. got=%s", fn.ReturnType)
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
