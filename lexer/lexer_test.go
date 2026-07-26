package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken_Basic(t *testing.T) {
	input := `var five = 5
ten := 10.5

func add(x, y) {
	return x + y
}

result := add(five, ten)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.IDENT, "ten"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.IDENT, "result"},
		{token.COLON_ASSIGN, ":="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_OperatorsAndCommentsAndStrings(t *testing.T) {
	input := `
// Single-line comment
x++
y--
/* Multi-line
   comment */
"hello\nworld"
/* comment1 */ // comment2
"adjacent"[`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "x"},
		{token.INC, "++"},
		{token.IDENT, "y"},
		{token.DEC, "--"},
		{token.STRING, "hello\nworld"},
		{token.STRING, "adjacent"},
		{token.LBRACKET, "["},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q at line %d col %d",
				i, tt.expectedType, tok.Type, tok.Line, tok.Column)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Positions(t *testing.T) {
	input := "x := 5\ny := 10"
	l := New(input)

	// x
	tok := l.NextToken()
	if tok.Line != 1 || tok.Column != 1 {
		t.Errorf("Expected 'x' at line 1 col 1, got line %d col %d", tok.Line, tok.Column)
	}

	// :=
	tok = l.NextToken()
	if tok.Line != 1 || tok.Column != 3 {
		t.Errorf("Expected ':=' at line 1 col 3, got line %d col %d", tok.Line, tok.Column)
	}

	// 5
	tok = l.NextToken()
	if tok.Line != 1 || tok.Column != 6 {
		t.Errorf("Expected '5' at line 1 col 6, got line %d col %d", tok.Line, tok.Column)
	}

	// y
	tok = l.NextToken()
	if tok.Line != 2 || tok.Column != 1 {
		t.Errorf("Expected 'y' at line 2 col 1, got line %d col %d", tok.Line, tok.Column)
	}
}
