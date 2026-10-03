package parser

import (
	"testing"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
)

func TestVarStatements(t *testing.T) {
	input := `
var x int = 5
var y = 10
var name string = "RustX"
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

func TestShorthandAssignmentAndIndexAssignment(t *testing.T) {
	input := `
x := 5
arr[0] = 10
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	assign1, ok1 := program.Statements[0].(*ast.AssignStatement)
	if !ok1 {
		t.Fatalf("stmt 0 not *ast.AssignStatement. got=%T", program.Statements[0])
	}
	ident, okIdent := assign1.Name.(*ast.Identifier)
	if !okIdent || ident.Value != "x" {
		t.Fatalf("assign1.Name not Identifier 'x'. got=%v", assign1.Name)
	}

	assign2, ok2 := program.Statements[1].(*ast.AssignStatement)
	if !ok2 {
		t.Fatalf("stmt 1 not *ast.AssignStatement. got=%T", program.Statements[1])
	}
	indexExp, okIndex := assign2.Name.(*ast.IndexExpression)
	if !okIndex {
		t.Fatalf("assign2.Name not *ast.IndexExpression. got=%T", assign2.Name)
	}
	if indexExp.Left.String() != "arr" {
		t.Fatalf("indexExp.Left not 'arr'. got=%s", indexExp.Left.String())
	}
}

func TestForLoopAndWhileLoop(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i)
}

for i < 10 {
	i++
}

while i > 0 {
	i--
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for1, ok1 := program.Statements[0].(*ast.ForStatement)
	if !ok1 {
		t.Fatalf("stmt 0 not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if for1.Init == nil || for1.Condition == nil || for1.Post == nil {
		t.Fatalf("for1 missing init, cond, or post")
	}

	for2, ok2 := program.Statements[1].(*ast.ForStatement)
	if !ok2 {
		t.Fatalf("stmt 1 not *ast.ForStatement. got=%T", program.Statements[1])
	}
	if for2.Init != nil || for2.Condition == nil {
		t.Fatalf("for2 expected single condition for loop")
	}

	while1, ok3 := program.Statements[2].(*ast.WhileStatement)
	if !ok3 {
		t.Fatalf("stmt 2 not *ast.WhileStatement. got=%T", program.Statements[2])
	}
	if while1.Condition == nil {
		t.Fatalf("while1 condition is nil")
	}
}

func TestFunctionLiteralAndCall(t *testing.T) {
	input := `
func add(a int, b int) int {
	return a + b
}

res := add(5, 10)
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	fnStmt, okFn := program.Statements[0].(*ast.ExpressionStatement)
	if !okFn {
		t.Fatalf("stmt 0 not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	fnLit, okLit := fnStmt.Expression.(*ast.FunctionLiteral)
	if !okLit {
		t.Fatalf("fnStmt not FunctionLiteral. got=%T", fnStmt.Expression)
	}
	if fnLit.Name != "add" {
		t.Errorf("fnLit.Name not 'add'. got=%s", fnLit.Name)
	}
	if len(fnLit.Parameters) != 2 {
		t.Fatalf("fnLit.Parameters length not 2. got=%d", len(fnLit.Parameters))
	}
	if fnLit.ReturnType != "int" {
		t.Errorf("fnLit.ReturnType not 'int'. got=%s", fnLit.ReturnType)
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
