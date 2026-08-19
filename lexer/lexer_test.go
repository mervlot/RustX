package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
var x = 5;
name := "RustX\nStarter";
func add(a int, b float) float {
	// single line comment
	/* multi
	   line
	   comment */
	if a == b && b != 0 {
		return a + b;
	} else if a <= b || b >= a {
		i++;
		j--;
	}
	return 0.0;
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "RustX\nStarter"},
		{token.SEMICOLON, ";"},

		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.TYPE_INT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.TYPE_FLT, "float"},
		{token.RPAREN, ")"},
		{token.TYPE_FLT, "float"},
		{token.LBRACE, "{"},

		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.LTE, "<="},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "b"},
		{token.GTE, ">="},
		{token.IDENT, "a"},
		{token.LBRACE, "{"},

		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "j"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.FLOAT, "0.0"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
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
