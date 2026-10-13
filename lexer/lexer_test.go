package lexer

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
var x = 10
y := 20.5
// single line comment
/* multi line
   comment */
func add(a int, b int) int {
	if a == b {
		return a + b
	} else if a != b {
		return a - b * 2 / 1 % 5
	}
	return 0
}

str := "Hello\nWorld"
for i := 0; i < 10; i++ {
	while i <= 5 && i >= 1 || i > 0 {
		println(str)
	}
}
arr := [1, 2, 3]
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "10"},

		{token.IDENT, "y"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "20.5"},

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
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},

		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.IDENT, "b"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.MINUS, "-"},
		{token.IDENT, "b"},
		{token.ASTERISK, "*"},
		{token.INT, "2"},
		{token.SLASH, "/"},
		{token.INT, "1"},
		{token.MOD, "%"},
		{token.INT, "5"},

		{token.RBRACE, "}"},

		{token.RETURN, "return"},
		{token.INT, "0"},

		{token.RBRACE, "}"},

		{token.IDENT, "str"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Hello\nWorld"},

		{token.FOR, "for"},
		{token.IDENT, "i"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.INC, "++"},
		{token.LBRACE, "{"},

		{token.WHILE, "while"},
		{token.IDENT, "i"},
		{token.LTE, "<="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "i"},
		{token.GTE, ">="},
		{token.INT, "1"},
		{token.OR, "||"},
		{token.IDENT, "i"},
		{token.GT, ">"},
		{token.INT, "0"},
		{token.LBRACE, "{"},

		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.IDENT, "str"},
		{token.RPAREN, ")"},

		{token.RBRACE, "}"},

		{token.RBRACE, "}"},

		{token.IDENT, "arr"},
		{token.COLON_ASSIGN, ":="},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.COMMA, ","},
		{token.INT, "3"},
		{token.RBRACKET, "]"},

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
