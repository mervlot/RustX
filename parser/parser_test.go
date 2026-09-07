package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5;
var y float = 10.5;
var name string = "RustX";
var z = 20;
a := "shorthand";
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 5 {
		t.Fatalf("program.Statements does not contain 5 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"name"},
		{"z"},
		{"a"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		switch s := stmt.(type) {
		case *ast.VarStatement:
			if s.Name.Value != tt.expectedIdentifier {
				t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, s.Name.Value)
			}
		case *ast.ShortVarStatement:
			if s.Name.Value != tt.expectedIdentifier {
				t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.expectedIdentifier, s.Name.Value)
			}
		default:
			t.Fatalf("stmt not *ast.VarStatement or *ast.ShortVarStatement. got=%T", stmt)
		}
	}
}

func TestFunctionDeclaration(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
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
		t.Fatalf("stmt is not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fn.Name.Value != "add" {
		t.Fatalf("function name not 'add'. got=%s", fn.Name.Value)
	}

	if len(fn.Parameters) != 2 {
		t.Fatalf("function parameters count wrong. got=%d", len(fn.Parameters))
	}

	if fn.ReturnType == nil || fn.ReturnType.Name != "int" {
		t.Fatalf("function return type wrong. got=%v", fn.ReturnType)
	}
}

func TestControlFlowStatements(t *testing.T) {
	input := `
if x < y {
	println("less")
} else {
	println("greater")
}

while x > 0 {
	x--
}

for i := 0; i < 10; i++ {
	println(i)
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
