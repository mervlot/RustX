package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var x = 5;
name := "RustX";
// single line comment
/* multiline
   comment */
func add(a int, b int) int {
	if (a == b) {
		return a + b;
	} else if (a != b && a >= 0) {
		a++;
		b--;
		return a * b / 2.5;
	}
	return 0;
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
		{token.STRING, "RustX"},
		{token.SEMICOLON, ";"},
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.IDENT, "b"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.GTE, ">="},
		{token.INT, "0"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "a"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "b"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.ASTERISK, "*"},
		{token.IDENT, "b"},
		{token.SLASH, "/"},
		{token.FLOAT, "2.5"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RETURN, "return"},
		{token.INT, "0"},
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
