package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var age = 16
name := "Muhammad"
// Single line comment
/* Multi-line
   comment */
func add(a int, b int) float {
	return a + b
}

if age >= 18 && name != "" || age == 16 {
	i++
	j--
	x % 2
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.INT, "16"},
		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Muhammad"},
		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "float"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.IDENT, "age"},
		{token.GTE, ">="},
		{token.INT, "18"},
		{token.AND, "&&"},
		{token.IDENT, "name"},
		{token.NOT_EQ, "!="},
		{token.STRING, ""},
		{token.OR, "||"},
		{token.IDENT, "age"},
		{token.EQ, "=="},
		{token.INT, "16"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.IDENT, "j"},
		{token.DEC, "--"},
		{token.IDENT, "x"},
		{token.MODULO, "%"},
		{token.INT, "2"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStringEscapes(t *testing.T) {
	input := `"hello\nworld\t\"quoted\"\\"`
	l := New(input)
	tok := l.NextToken()

	if tok.Type != token.STRING {
		t.Fatalf("tok.Type not STRING. got=%q", tok.Type)
	}
	expected := "hello\nworld\t\"quoted\"\\"
	if tok.Literal != expected {
		t.Fatalf("tok.Literal wrong. expected=%q, got=%q", expected, tok.Literal)
	}
}
