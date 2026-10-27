package parser

import (
	"testing"

	"rustx/ast"
	"rustx/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

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

func TestShorthandAndAssign(t *testing.T) {
	input := `
x := 5;
x = 10;
arr[0] = 42;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	assign1, ok := program.Statements[0].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt 0 not *ast.AssignStatement. got=%T", program.Statements[0])
	}
	if assign1.Token.Literal != ":=" {
		t.Errorf("assign1 token literal expected ':=', got=%s", assign1.Token.Literal)
	}

	assign2, ok := program.Statements[1].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt 1 not *ast.AssignStatement. got=%T", program.Statements[1])
	}
	if assign2.Token.Literal != "=" {
		t.Errorf("assign2 token literal expected '=', got=%s", assign2.Token.Literal)
	}

	assign3, ok := program.Statements[2].(*ast.AssignStatement)
	if !ok {
		t.Fatalf("stmt 2 not *ast.AssignStatement. got=%T", program.Statements[2])
	}
	_, ok = assign3.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("assign3.Left not *ast.IndexExpression. got=%T", assign3.Left)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `func add(a int, b int) int { return a + b; }`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	fnStmt, ok := program.Statements[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt not *ast.FunctionLiteral. got=%T", program.Statements[0])
	}

	if fnStmt.Name != "add" {
		t.Errorf("function name expected 'add', got=%s", fnStmt.Name)
	}
	if len(fnStmt.Parameters) != 2 {
		t.Fatalf("function parameters count expected 2, got=%d", len(fnStmt.Parameters))
	}
	if fnStmt.Parameters[0].Value != "a" || fnStmt.ParamTypes[0] != "int" {
		t.Errorf("param 0 wrong: %s %s", fnStmt.Parameters[0].Value, fnStmt.ParamTypes[0])
	}
	if fnStmt.ReturnType != "int" {
		t.Errorf("return type expected 'int', got=%s", fnStmt.ReturnType)
	}
}

func TestForAndWhileLoops(t *testing.T) {
	input := `
for i := 0; i < 10; i++ {
	println(i);
}

for i < 10 {
	i++;
}

while x > 0 {
	x--;
}
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	forStmt, ok := program.Statements[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt 0 not *ast.ForStatement. got=%T", program.Statements[0])
	}
	if forStmt.Init == nil || forStmt.Condition == nil || forStmt.Post == nil {
		t.Errorf("three-part for loop missing components: init=%v, cond=%v, post=%v", forStmt.Init, forStmt.Condition, forStmt.Post)
	}

	singleFor, ok := program.Statements[1].(*ast.ForStatement)
	if !ok {
		t.Fatalf("stmt 1 not *ast.ForStatement. got=%T", program.Statements[1])
	}
	if singleFor.Condition == nil || singleFor.Init != nil {
		t.Errorf("single condition for loop parsed incorrectly")
	}

	whileStmt, ok := program.Statements[2].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("stmt 2 not *ast.WhileStatement. got=%T", program.Statements[2])
	}
	if whileStmt.Condition == nil {
		t.Errorf("while statement condition is nil")
	}
}

func TestArrayLiteralAndIndex(t *testing.T) {
	input := `arr := [1, 2, 3]; x := arr[1];`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d", len(program.Statements))
	}

	assign1 := program.Statements[0].(*ast.AssignStatement)
	_, ok := assign1.Value.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("assign1 value not *ast.ArrayLiteral. got=%T", assign1.Value)
	}

	assign2 := program.Statements[1].(*ast.AssignStatement)
	_, ok = assign2.Value.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("assign2 value not *ast.IndexExpression. got=%T", assign2.Value)
	}
}
