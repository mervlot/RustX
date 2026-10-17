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
var foobar float = 8.3;
`
	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("unexpected parser init error: %v", err)
	}

	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram() failed: %v", err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
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
	}
}

func TestShortVarStatements(t *testing.T) {
	input := `x := 5;`
	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("unexpected parser init error: %v", err)
	}

	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram() failed: %v", err)
	}

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ShortVarStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ShortVarStatement. got=%T", program.Statements[0])
	}

	if stmt.Name.Value != "x" {
		t.Errorf("stmt.Name.Value not 'x'. got=%s", stmt.Name.Value)
	}
}

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b;
}

result := add(5, 10);
`
	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("unexpected parser init error: %v", err)
	}

	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram() failed: %v", err)
	}

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fnStmt, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt[0] not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fnStmt.Name.Value != "add" {
		t.Errorf("fnStmt.Name not 'add'. got=%s", fnStmt.Name.Value)
	}
	if len(fnStmt.Parameters) != 2 {
		t.Errorf("fnStmt parameters count wrong. got=%d", len(fnStmt.Parameters))
	}
}

func TestControlFlowParsing(t *testing.T) {
	input := `
if x < y {
	println("smaller");
} else {
	println("larger");
}

while x < 10 {
	x++;
}

for i := 0; i < 10; i++ {
	println(i);
}
`
	l := lexer.New(input)
	p, err := New(l)
	if err != nil {
		t.Fatalf("unexpected parser init error: %v", err)
	}

	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("ParseProgram() failed: %v", err)
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
}
