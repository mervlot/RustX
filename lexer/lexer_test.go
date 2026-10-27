package lexer

import (
	"testing"

	"rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var x = 10;
// single line comment
/* multiline
   comment */
name := "Hello\nWorld";
func add(a int, b float) int {
	if a <= 5 && b >= 3.14 {
		return a + 1;
	} else if a != 0 || b == 0.0 {
		a++;
		b--;
		a % 2;
	}
	return b;
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Hello\nWorld"},
		{token.SEMICOLON, ";"},

		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.IDENT, "float"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},

		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.LT_EQ, "<="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.GT_EQ, ">="},
		{token.FLOAT, "3.14"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.INT, "0"},
		{token.OR, "||"},
		{token.IDENT, "b"},
		{token.EQ, "=="},
		{token.FLOAT, "0.0"},
		{token.LBRACE, "{"},

		{token.IDENT, "a"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "b"},
		{token.DEC, "--"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "a"},
		{token.MOD, "%"},
		{token.INT, "2"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},

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
